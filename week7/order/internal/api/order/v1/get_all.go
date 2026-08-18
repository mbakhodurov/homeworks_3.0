package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week7/order/internal/converter"
	orderv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/openapi/order/v1"
)

// GetAllOrders возвращает список всех заказов.
func (a *api) GetAllOrders(ctx context.Context) (orderv1.GetAllOrdersRes, error) {
	orders, err := a.orderService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]orderv1.OrderDto, 0, len(orders))
	for _, o := range orders {
		dtos = append(dtos, *converter.OrderToDTO(o))
	}

	return &orderv1.GetAllOrdersResponse{
		Orders:     dtos,
		TotalCount: int64(len(dtos)),
	}, nil
}
