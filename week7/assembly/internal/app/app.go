package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbakhodurov/homeworks2/week7/assembly/internal/config"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	diContainer *diContainer
}

func New() *App {
	a := &App{}
	a.initDeps()
	return a
}

func (a *App) initDeps() {
	inits := []func(){
		a.initDI,
		a.initLogger,
	}

	for _, f := range inits {
		f()
	}
}

func (a *App) initDI() {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger() {
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

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runConsumer(ctx)
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

func (a *App) runConsumer(ctx context.Context) error {
	slog.Info("Kafka-потребитель OrderPaid запущен")

	return a.diContainer.OrderPaidConsumer().Run(ctx)
}
