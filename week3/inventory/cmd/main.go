package main

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
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	inventoryapi "github.com/mbakhodurov/homeworks2/week3/inventory/internal/api/inventory/v1"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/interceptor"
	partrepo "github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/part"
	partsvc "github.com/mbakhodurov/homeworks2/week3/inventory/internal/service/part"
	inventory_v1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/inventory/v1"
	"github.com/mbakhodurov/homeworks2/week3/shared/static"
)

const (
	inventoryDSN = "postgres://inventory-service-user:inventory-service-password@localhost:5433/inventory-service?sslmode=disable"

	grpcAddress = "localhost:50051"
	httpAddress = ":8083"

	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute

	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 15 * time.Second
	httpWriteTimeout      = 15 * time.Second
	httpIdleTimeout       = 60 * time.Second
	httpShutdownTimeout   = 5 * time.Second

	swaggerUIFile   = "swagger-ui.html"
	swaggerJSONFile = "generated/inventory/v1/inventory.swagger.json"
)

//nolint:funlen
func main() {
	bgCtx := context.Background()

	pool, err := pgxpool.New(bgCtx, inventoryDSN)
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(bgCtx); err != nil {
		slog.Error("проверка соединения с БД", "error", err)
		os.Exit(1)
	}
	slog.Info("подключение к PostgreSQL установлено")

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor(),
			interceptor.LoggerInterceptor(),
			interceptor.ErrorInterceptor,
		),
	)

	repo := partrepo.NewRepository(pool)
	svc := partsvc.New(repo)
	api := inventoryapi.New(svc)

	inventory_v1.RegisterInventoryServiceServer(grpcServer, api)
	reflection.Register(grpcServer)

	ctx, cancel := signal.NotifyContext(bgCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("запуск InventoryService gRPC", "адрес", grpcAddress)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			slog.Error("ошибка запуска gRPC сервера", "error", serveErr)
			cancel()
		}
	}()

	gwCtx, gwCancel := context.WithCancel(bgCtx)
	defer gwCancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err = inventory_v1.RegisterInventoryServiceHandlerFromEndpoint(gwCtx, mux, grpcAddress, opts); err != nil {
		slog.Error("ошибка регистрации gateway", "error", err)
		gwCancel()
		return
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", mux)

	httpMux.Handle("/swagger-ui.html", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, readErr := static.FS.ReadFile(swaggerUIFile)
		if readErr != nil {
			http.Error(w, "swagger-ui.html не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, writeErr := w.Write(data); writeErr != nil {
			slog.Error("ошибка записи swagger-ui", "error", writeErr)
		}
	}))

	httpMux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, readErr := static.FS.ReadFile(swaggerJSONFile)
		if readErr != nil {
			http.Error(w, "swagger.json не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, writeErr := w.Write(data); writeErr != nil {
			slog.Error("ошибка записи swagger.json", "error", writeErr)
		}
	}))

	httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	}))

	gwServer := &http.Server{
		Addr:              httpAddress,
		Handler:           httpMux,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	go func() {
		slog.Info("запуск InventoryService HTTP + Swagger UI", "адрес", httpAddress)
		if serveErr := gwServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска HTTP сервера", "error", serveErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("завершение работы серверов...")

	shutdownCtx, shutdownCancel := context.WithTimeout(bgCtx, httpShutdownTimeout)
	defer shutdownCancel()
	if shutdownErr := gwServer.Shutdown(shutdownCtx); shutdownErr != nil {
		slog.Error("ошибка остановки HTTP сервера", "error", shutdownErr)
	}
	slog.Info("HTTP сервер остановлен")

	grpcServer.GracefulStop()
	slog.Info("gRPC сервер остановлен")
}
