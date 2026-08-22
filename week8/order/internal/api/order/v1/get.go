package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/order/internal/converter"
	orderv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/openapi/order/v1"
)

// GetOrder возвращает полную информацию о заказе.
func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return converter.OrderToDTO(order), nil
}
