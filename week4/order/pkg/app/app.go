package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	orderapiv1 "github.com/mbakhodurov/homeworks2/week4/order/internal/api/order/v1"
	inventoryclient "github.com/mbakhodurov/homeworks2/week4/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/mbakhodurov/homeworks2/week4/order/internal/client/grpc/payment/v1"
	orderrepo "github.com/mbakhodurov/homeworks2/week4/order/internal/repository/order"
	ordersvc "github.com/mbakhodurov/homeworks2/week4/order/internal/service/order"
	orderv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler собирает зависимости и возвращает HTTP-обработчик OrderService.
// Используется в API-тестах: тестовый код в модуле order не может импортировать
// internal/-пакеты из модулей inventory и payment, поэтому клиенты этих сервисов
// передаются снаружи уже как готовые gRPC-стабы (например, поднятые через bufconn).
func NewHTTPHandler(
	pool *pgxpool.Pool,
	txManager ordersvc.TxManager,
	invClient inventoryv1.InventoryServiceClient,
	payClient paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	orderRepo := orderrepo.NewRepository(pool)
	inventoryClient := inventoryclient.New(invClient)
	paymentClient := paymentclient.New(payClient)

	svc := ordersvc.New(orderRepo, inventoryClient, paymentClient, txManager)
	api := orderapiv1.New(svc)

	return orderv1.NewServer(api, orderv1.WithErrorHandler(orderapiv1.ErrorHandler))
}
