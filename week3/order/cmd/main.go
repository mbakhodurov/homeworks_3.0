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

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderapi "github.com/mbakhodurov/homeworks2/week3/order/internal/api/order/v1"
	inventoryclient "github.com/mbakhodurov/homeworks2/week3/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/mbakhodurov/homeworks2/week3/order/internal/client/grpc/payment/v1"
	orderrepo "github.com/mbakhodurov/homeworks2/week3/order/internal/repository/order"
	orderitemrepo "github.com/mbakhodurov/homeworks2/week3/order/internal/repository/order_item"
	ordersvc "github.com/mbakhodurov/homeworks2/week3/order/internal/service/order"
	orderv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/payment/v1"
	"github.com/mbakhodurov/homeworks2/week3/shared/static"
)

const (
	httpPort = "8080"

	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"

	swaggerUIFile    = "swagger-ui.html"
	swaggerOrderFile = "generated/order.swagger.yaml"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	bgCtx := context.Background()

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		dbURI = "postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable"
	}

	pool, err := pgxpool.New(bgCtx, dbURI)
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

	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return
	}
	defer paymentConn.Close()

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("создание transaction manager", "error", err)
		os.Exit(1)
	}

	repo := orderrepo.NewRepository(pool)
	orderItemRepo := orderitemrepo.NewRepository(pool)
	invClient := inventoryclient.New(inventoryv1.NewInventoryServiceClient(inventoryConn))
	payClient := paymentclient.New(paymentv1.NewPaymentServiceClient(paymentConn))
	svc := ordersvc.New(repo, orderItemRepo, invClient, payClient, txManager)
	api := orderapi.New(svc)

	orderServer, err := orderv1.NewServer(api, orderv1.WithErrorHandler(orderapi.ErrorHandler))
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/api/", orderServer)

	mux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, _ *http.Request) {
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

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		data, readErr := static.FS.ReadFile(swaggerOrderFile)
		if readErr != nil {
			http.Error(w, "swagger.json не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		if _, writeErr := w.Write(data); writeErr != nil {
			slog.Error("ошибка записи swagger.json", "error", writeErr)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("", httpPort),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, cancel := signal.NotifyContext(bgCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("запуск OrderService", "port", httpPort)
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("завершение работы сервера...")

	shutdownCtx, cancelShutdown := context.WithTimeout(bgCtx, shutdownTimeout)
	defer cancelShutdown()

	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		slog.Error("ошибка при остановке сервера", "error", shutdownErr)
	}

	slog.Info("сервер остановлен")
}
