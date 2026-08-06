package app

import (
	"context"
	"errors"
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
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/config"
	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/interceptor"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/logger"
	inventory_v1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/inventory/v1"
	"github.com/mbakhodurov/homeworks2/week6/shared/static"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 5 * time.Second

	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 15 * time.Second
	httpWriteTimeout      = 15 * time.Second
	httpIdleTimeout       = 60 * time.Second

	swaggerUIFile   = "swagger-ui.html"
	swaggerJSONFile = "generated/inventory/v1/inventory.swagger.json"
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
		grpc.Creds(insecure.NewCredentials()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: false,
		}),
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor(),
			interceptor.ErrorInterceptor,
			interceptor.LoggerInterceptor(),
			interceptor.Auth(a.diContainer.IAMClient()),
		),
	)

	api := a.diContainer.PartAPI(ctx)

	closer.Add("gRPC сервер", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	inventory_v1.RegisterInventoryServiceServer(a.grpcServer, api)
}

func (a *App) initHTTPServer(_ context.Context) {
	gwCtx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := inventory_v1.RegisterInventoryServiceHandlerFromEndpoint(gwCtx, mux, config.AppConfig().GRPC.Address(), opts); err != nil {
		slog.Error("ошибка регистрации grpc-gateway", "error", err)
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
		data, readErr := static.FS.ReadFile(swaggerJSONFile)
		if readErr != nil {
			http.Error(w, "swagger.json не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write(data); writeErr != nil {
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
	slog.Info("запуск InventoryService gRPC", "адрес", config.AppConfig().GRPC.Address())
	return a.grpcServer.Serve(a.listener)
}

func (a *App) runHTTPServer() error {
	slog.Info("запуск InventoryService HTTP + Swagger UI", "адрес", config.AppConfig().HTTP.Address())
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
