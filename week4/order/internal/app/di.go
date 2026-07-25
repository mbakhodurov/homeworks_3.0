package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderapiv1 "github.com/mbakhodurov/homeworks2/week4/order/internal/api/order/v1"
	inventoryclient "github.com/mbakhodurov/homeworks2/week4/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/mbakhodurov/homeworks2/week4/order/internal/client/grpc/payment/v1"
	"github.com/mbakhodurov/homeworks2/week4/order/internal/config"
	orderrepo "github.com/mbakhodurov/homeworks2/week4/order/internal/repository/order"
	orderitemrepo "github.com/mbakhodurov/homeworks2/week4/order/internal/repository/order_item"
	ordersvc "github.com/mbakhodurov/homeworks2/week4/order/internal/service/order"
	"github.com/mbakhodurov/homeworks2/week4/platform/pkg/closer"
	orderv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/payment/v1"
)

// diContainer — контейнер зависимостей с ленивой инициализацией.
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует результат.
type diContainer struct {
	pgPool    *pgxpool.Pool
	txManager ordersvc.TxManager

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	orderRepo       ordersvc.OrderRepository
	orderItemRepo   ordersvc.OrderItemRepository
	inventoryClient ordersvc.InventoryClient
	paymentClient   ordersvc.PaymentClient

	orderService orderapiv1.OrderService

	httpHandler http.Handler
}

// PGPool возвращает пул подключений к PostgreSQL.
// При первом вызове создаёт пул, проверяет соединение и регистрирует closer.
func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		if err = pool.Ping(ctx); err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

// TxManager возвращает менеджер транзакций.
func (d *diContainer) TxManager(ctx context.Context) ordersvc.TxManager {
	if d.txManager == nil {
		m, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = m
	}

	return d.txManager
}

// InventoryConn возвращает gRPC-соединение с InventoryService.
func (d *diContainer) InventoryConn() *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			os.Exit(1)
		}

		closer.Add("InventoryService gRPC-соединение", func(_ context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

// PaymentConn возвращает gRPC-соединение с PaymentService.
func (d *diContainer) PaymentConn() *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			os.Exit(1)
		}

		closer.Add("PaymentService gRPC-соединение", func(_ context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

// OrderRepository возвращает репозиторий заказов.
func (d *diContainer) OrderRepository(ctx context.Context) ordersvc.OrderRepository {
	if d.orderRepo == nil {
		d.orderRepo = orderrepo.NewRepository(d.PGPool(ctx))
	}

	return d.orderRepo
}

// OrderItemRepository возвращает репозиторий позиций заказа.
func (d *diContainer) OrderItemRepository(ctx context.Context) ordersvc.OrderItemRepository {
	if d.orderItemRepo == nil {
		d.orderItemRepo = orderitemrepo.NewRepository(d.PGPool(ctx))
	}

	return d.orderItemRepo
}

// InventoryClient возвращает клиент InventoryService.
func (d *diContainer) InventoryClient() ordersvc.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryclient.New(inventoryv1.NewInventoryServiceClient(d.InventoryConn()))
	}

	return d.inventoryClient
}

// PaymentClient возвращает клиент PaymentService.
func (d *diContainer) PaymentClient() ordersvc.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentclient.New(paymentv1.NewPaymentServiceClient(d.PaymentConn()))
	}

	return d.paymentClient
}

// OrderService возвращает application-сервис заказов.
func (d *diContainer) OrderService(ctx context.Context) orderapiv1.OrderService {
	if d.orderService == nil {
		d.orderService = ordersvc.New(
			d.OrderRepository(ctx),
			d.OrderItemRepository(ctx),
			d.InventoryClient(),
			d.PaymentClient(),
			d.TxManager(ctx),
		)
	}

	return d.orderService
}

// HTTPHandler возвращает HTTP-обработчик OrderService (ogen-сервер поверх OpenAPI-спеки).
func (d *diContainer) HTTPHandler(ctx context.Context) http.Handler {
	if d.httpHandler == nil {
		api := orderapiv1.New(d.OrderService(ctx))

		server, err := orderv1.NewServer(api, orderv1.WithErrorHandler(orderapiv1.ErrorHandler))
		if err != nil {
			slog.Error("не удалось создать HTTP-сервер OpenAPI", "error", err)
			os.Exit(1)
		}

		d.httpHandler = server
	}

	return d.httpHandler
}
