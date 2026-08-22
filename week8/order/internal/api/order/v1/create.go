package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/order/internal/converter"
	orderv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/openapi/order/v1"
)

// CreateOrder создаёт новый заказ на постройку космического корабля.
func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	order, err := a.orderService.Create(ctx, converter.ToCreateOrderInput(req))
	if err != nil {
		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}
