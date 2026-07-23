package app

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapiv1 "github.com/mbakhodurov/homeworks2/week4/inventory/internal/api/inventory/v1"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/interceptor"
	partrepo "github.com/mbakhodurov/homeworks2/week4/inventory/internal/repository/part"
	partsvc "github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/application/part"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/domain"
	inventory_v1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// Interceptors возвращает gRPC-опции с интерцепторами для использования в тестах.
func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor(),
			interceptor.ErrorInterceptor,
			interceptor.LoggerInterceptor(),
		),
	}
}

// RegisterServices собирает зависимости и регистрирует gRPC-сервисы на сервере.
func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		panic(err)
	}

	partRepo := partrepo.NewRepository(pool)
	checker := domain.NewCompatibilityChecker()

	svc := partsvc.New(txManager, partRepo, checker)
	api := inventoryapiv1.New(svc)

	inventory_v1.RegisterInventoryServiceServer(grpcServer, api)
}
