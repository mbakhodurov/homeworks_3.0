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

	"github.com/mbakhodurov/homeworks2/week5/payment/internal/config"
	"github.com/mbakhodurov/homeworks2/week5/payment/internal/interceptor"
	"github.com/mbakhodurov/homeworks2/week5/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week5/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks2/week5/platform/pkg/logger"
	payment_v1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/proto/payment/v1"
	"github.com/mbakhodurov/homeworks2/week5/shared/static"
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
	swaggerJSONFile = "generated/payment/v1/payment.swagger.json"
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
	httpServer  *http.Server
}

// New создаёт и инициализирует приложение.
func New() *App {
	a := &App{}
	a.initDeps()
	return a
}

// initDeps последовательно инициализирует все зависимости приложения.
func (a *App) initDeps() {
	inits := []func(){
		a.initDI,
		a.initLogger,
		a.initListener,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		f()
	}
}

func (a *App) initDI() {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger() {
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) initListener() {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address()) //nolint:noctx
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}

	a.listener = listener
}

func (a *App) initGRPCServer() {
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
		),
	)

	api := a.diContainer.PaymentAPI()

	closer.Add("gRPC сервер", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	payment_v1.RegisterPaymentServiceServer(a.grpcServer, api)
}

func (a *App) initHTTPServer() {
	gwCtx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := payment_v1.RegisterPaymentServiceHandlerFromEndpoint(gwCtx, mux, config.AppConfig().GRPC.Address(), opts); err != nil {
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
	slog.Info("запуск PaymentService gRPC", "адрес", config.AppConfig().GRPC.Address())
	return a.grpcServer.Serve(a.listener)
}

func (a *App) runHTTPServer() error {
	slog.Info("запуск PaymentService HTTP + Swagger UI", "адрес", config.AppConfig().HTTP.Address())
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
