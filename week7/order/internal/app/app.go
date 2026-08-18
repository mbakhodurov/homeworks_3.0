package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mbakhodurov/homeworks2/week7/order/internal/config"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/logger"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/metrics"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/tracing"
)

const (
	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 15 * time.Second
	httpWriteTimeout      = 15 * time.Second
	httpIdleTimeout       = 60 * time.Second
	httpShutdownTimeout   = 5 * time.Second
	shutdownTimeout       = 5 * time.Second
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов.
type App struct {
	diContainer *diContainer
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
		a.initTracer,
		a.initMetrics,
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
	cfg := config.AppConfig()
	logger.Init(logger.Config{
		Level:             cfg.Logger.Level,
		ServiceName:       cfg.OTel.ServiceName,
		Environment:       "development",
		EnableOTLP:        true,
		CollectorEndpoint: cfg.OTel.Endpoint,
	})
	closer.Add("logger", func(_ context.Context) error {
		return logger.Close()
	})
}

func (a *App) initTracer(ctx context.Context) {
	cfg := config.AppConfig()
	shutdown, err := tracing.InitTracer(ctx, tracing.Config{
		CollectorEndpoint: cfg.OTel.Endpoint,
		ServiceName:       cfg.OTel.ServiceName,
		Environment:       "development",
	})
	if err != nil {
		slog.Error("не удалось инициализировать трейсер", "error", err)
		os.Exit(1)
	}
	closer.Add("tracer", shutdown)
}

func (a *App) initMetrics(_ context.Context) {
	cfg := config.AppConfig()
	metrics.Init(cfg.OTel.ServiceName)
	closer.Add("metrics", func(_ context.Context) error {
		return metrics.Close()
	})
}

func (a *App) initHTTPServer(ctx context.Context) {
	handler := a.diContainer.HTTPHandler(ctx)

	// otelhttp автоматически создаёт span для каждого HTTP-запроса
	// и добавляет метрику http_server_request_duration_seconds
	instrumentedHandler := otelhttp.NewHandler(handler, "order-service")

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           instrumentedHandler,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	closer.Add("HTTP сервер", func(shutdownCtx context.Context) error {
		return a.httpServer.Shutdown(shutdownCtx)
	})
}

// Run запускает HTTP-сервер и управляет жизненным циклом приложения.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runHTTPServer()
	}()

	go func() {
		if err := a.diContainer.AssemblyConsumer(ctx).Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			slog.ErrorContext(ctx, "consumer ShipAssembled завершился с ошибкой", "error", err)
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

func (a *App) runHTTPServer() error {
	slog.Info("запуск OrderService HTTP", "адрес", config.AppConfig().HTTP.Address())

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
