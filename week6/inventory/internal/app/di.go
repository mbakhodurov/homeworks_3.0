package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryapiv1 "github.com/mbakhodurov/homeworks2/week6/inventory/internal/api/inventory/v1"
	iamclient "github.com/mbakhodurov/homeworks2/week6/inventory/internal/client/grpc/iam/v1"
	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/config"
	partrepo "github.com/mbakhodurov/homeworks2/week6/inventory/internal/repository/part"
	partsvc "github.com/mbakhodurov/homeworks2/week6/inventory/internal/service/application/part"
	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/service/domain"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/closer"
	authv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/inventory/v1"
)

// diContainer — контейнер зависимостей с ленивой инициализацией.
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует результат.
type diContainer struct {
	pgPool    *pgxpool.Pool
	txManager partsvc.TxManager

	partRepo partsvc.PartRepository
	checker  partsvc.CompatibilityChecker

	partService inventoryapiv1.PartService

	partHandler inventory_v1.InventoryServiceServer

	iamConn   *grpc.ClientConn
	iamClient *iamclient.Client
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
func (d *diContainer) TxManager(ctx context.Context) partsvc.TxManager {
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

// PartRepository возвращает репозиторий деталей.
func (d *diContainer) PartRepository(ctx context.Context) partsvc.PartRepository {
	if d.partRepo == nil {
		d.partRepo = partrepo.NewRepository(d.PGPool(ctx))
	}

	return d.partRepo
}

// CompatibilityChecker возвращает доменный сервис проверки совместимости.
func (d *diContainer) CompatibilityChecker() partsvc.CompatibilityChecker {
	if d.checker == nil {
		d.checker = domain.NewCompatibilityChecker()
	}

	return d.checker
}

// PartService возвращает application-сервис деталей.
func (d *diContainer) PartService(ctx context.Context) inventoryapiv1.PartService {
	if d.partService == nil {
		d.partService = partsvc.New(
			d.TxManager(ctx),
			d.PartRepository(ctx),
			d.CompatibilityChecker(),
		)
	}

	return d.partService
}

// PartAPI возвращает gRPC-обработчик InventoryService.
func (d *diContainer) PartAPI(ctx context.Context) inventory_v1.InventoryServiceServer {
	if d.partHandler == nil {
		d.partHandler = inventoryapiv1.New(d.PartService(ctx))
	}

	return d.partHandler
}

// IAMConn возвращает gRPC-соединение с IAMService.
func (d *diContainer) IAMConn() *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к IAMService", "error", err)
			os.Exit(1)
		}

		closer.Add("IAMService gRPC-соединение", func(_ context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn
}

// IAMClient возвращает клиент AuthService (IAMService), нужный auth-interceptor'у.
func (d *diContainer) IAMClient() *iamclient.Client {
	if d.iamClient == nil {
		d.iamClient = iamclient.New(authv1.NewAuthServiceClient(d.IAMConn()))
	}

	return d.iamClient
}
