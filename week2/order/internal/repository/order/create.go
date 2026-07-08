package order

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/repository/converter"
)

// Create сохраняет новый заказ.
func (r *repository) Create(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.UUID] = converter.OrderToRecord(order)

	return nil
}
