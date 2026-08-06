package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/config"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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

// initGRPCServer строит gRPC-сервер через NewGRPCServer — ту же точку сборки,
// что использует API-тест iam/tests (тот же модуль — доступ к internal/ разрешён).
func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = NewGRPCServer(
		a.diContainer.PGPool(ctx),
		a.diContainer.RedisClient(ctx),
		config.AppConfig().Session.TTL,
		bcrypt.DefaultCost,
	)

	closer.Add("gRPC сервер", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})
}

// Run запускает gRPC-сервер и управляет жизненным циклом приложения.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runGRPCServer()
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
