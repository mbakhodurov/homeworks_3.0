package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/config"
	"github.com/mbakhodurov/homeworks2/week6/iam/internal/interceptor"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/logger"
	authv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/auth/v1"
	userv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/user/v1"
	"github.com/mbakhodurov/homeworks2/week6/shared/static"
)

const (
	shutdownTimeout = 5 * time.Second

	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 15 * time.Second
	httpWriteTimeout      = 15 * time.Second
	httpIdleTimeout       = 60 * time.Second

	swaggerUIFile       = "swagger-ui.html"
	swaggerAuthJSONFile = "generated/auth/v1/auth.swagger.json"
	swaggerUserJSONFile = "generated/user/v1/user.swagger.json"
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
	httpServer  *http.Server
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context) *App {
	a := &App{}
	a.initDeps(ctx)
	return a
}

// initDeps последовательно инициализирует все зависимости приложения.
func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initListener,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) initListener(_ context.Context) {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address()) //nolint:noctx // net.Listen не принимает context
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}

	a.listener = listener
}

func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.ErrorInterceptor),
	)

	authAPI := a.diContainer.AuthAPI(ctx)
	userAPI := a.diContainer.UserAPI(ctx)

	closer.Add("gRPC сервер", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	authv1.RegisterAuthServiceServer(a.grpcServer, authAPI)
	userv1.RegisterUserServiceServer(a.grpcServer, userAPI)
}

func (a *App) initHTTPServer(_ context.Context) {
	gwCtx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(gwCtx, mux, config.AppConfig().GRPC.Address(), opts); err != nil {
		slog.Error("ошибка регистрации grpc-gateway для AuthService", "error", err)
		os.Exit(1)
	}

	if err := userv1.RegisterUserServiceHandlerFromEndpoint(gwCtx, mux, config.AppConfig().GRPC.Address(), opts); err != nil {
		slog.Error("ошибка регистрации grpc-gateway для UserService", "error", err)
		os.Exit(1)
	}

	swaggerJSON, err := buildSwaggerJSON()
	if err != nil {
		slog.Error("не удалось собрать swagger.json", "error", err)
		os.Exit(1)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", mux)

	httpMux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, _ *http.Request) {
		data, readErr := static.FS.ReadFile(swaggerUIFile)
		if readErr != nil {
			http.Error(w, "swagger-ui.html не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, writeErr := w.Write(data); writeErr != nil {
			slog.Error("ошибка записи swagger-ui", "error", writeErr)
		}
	})

	httpMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write(swaggerJSON); writeErr != nil {
			slog.Error("ошибка записи swagger.json", "error", writeErr)
		}
	})

	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           httpMux,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	closer.Add("HTTP сервер", func(shutdownCtx context.Context) error {
		return a.httpServer.Shutdown(shutdownCtx)
	})
}

// buildSwaggerJSON объединяет сгенерированные swagger-документы AuthService и
// UserService в один: в отличие от inventory/payment, где на сервис приходится
// один proto-файл и swagger.json отдаётся статически как есть, у IAM их два.
func buildSwaggerJSON() ([]byte, error) {
	authDoc, err := readSwaggerDoc(swaggerAuthJSONFile)
	if err != nil {
		return nil, err
	}

	userDoc, err := readSwaggerDoc(swaggerUserJSONFile)
	if err != nil {
		return nil, err
	}

	merged := authDoc
	merged["info"] = map[string]any{
		"title":       "IAM Service",
		"description": "Package iam содержит сервисы аутентификации и управления пользователями.",
		"version":     "version not set",
	}
	merged["tags"] = append(toSlice(authDoc["tags"]), toSlice(userDoc["tags"])...)

	mergeInto(merged, userDoc, "paths")
	mergeInto(merged, userDoc, "definitions")

	swaggerJSON, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("сериализовать swagger.json: %w", err)
	}

	return swaggerJSON, nil
}

func readSwaggerDoc(file string) (map[string]any, error) {
	data, err := static.FS.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("прочитать %s: %w", file, err)
	}

	var doc map[string]any
	if err = json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("распарсить %s: %w", file, err)
	}

	return doc, nil
}

// mergeInto копирует ключи dst[key]/src[key] (оба — object в swagger-документе)
// из src в dst, создавая dst[key] при необходимости.
func mergeInto(dst, src map[string]any, key string) {
	dstMap, _ := dst[key].(map[string]any)
	if dstMap == nil {
		dstMap = map[string]any{}
	}

	srcMap, _ := src[key].(map[string]any)
	for k, v := range srcMap {
		dstMap[k] = v
	}

	dst[key] = dstMap
}

func toSlice(v any) []any {
	s, _ := v.([]any)
	return s
}

// Run запускает gRPC и HTTP серверы и управляет жизненным циклом приложения.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runGRPCServer()
	}()

	go func() {
		if err := a.runHTTPServer(); err != nil {
			slog.Error("HTTP сервер завершился с ошибкой", "error", err)
			cancel()
		}
	}()

	var runErr error
	select {
	case runErr = <-errCh:
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка при завершении работы", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

func (a *App) runGRPCServer() error {
	slog.Info("запуск IAMService gRPC", "адрес", config.AppConfig().GRPC.Address())
	return a.grpcServer.Serve(a.listener)
}

func (a *App) runHTTPServer() error {
	slog.Info("запуск IAMService HTTP + Swagger UI", "адрес", config.AppConfig().HTTP.Address())
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
