package v1

import (
	orderv1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/openapi/order/v1"
)

type api struct {
	orderv1.UnimplementedHandler

	orderService OrderService
}

// New создаёт новый HTTP API OrderService.
func New(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
