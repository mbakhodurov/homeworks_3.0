package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/api/converter"
	orderv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/openapi/order/v1"
)

// GetAllOrders возвращает список всех заказов.
func (a *api) GetAllOrders(ctx context.Context) (orderv1.GetAllOrdersRes, error) {
	orders, err := a.orderService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dto := converter.OrdersToDTO(orders)

	return &orderv1.GetAllOrdersResponse{Orders: dto, TotalCount: int64(len(dto))}, nil
}
