package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/mbakhodurov/homeworks2/week4/payment/internal/config"
	"github.com/mbakhodurov/homeworks2/week4/payment/internal/interceptor"
	"github.com/mbakhodurov/homeworks2/week4/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks2/week4/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks2/week4/platform/pkg/logger"
	payment_v1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/payment/v1"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 5 * time.Second
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
	slog.Info("запуск PaymentService gRPC", "адрес", config.AppConfig().GRPC.Address())
	return a.grpcServer.Serve(a.listener)
}
