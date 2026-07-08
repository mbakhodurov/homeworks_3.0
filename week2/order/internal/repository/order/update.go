package order

import (
	"context"

	errs "github.com/mbakhodurov/homeworks2/week2/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/repository/converter"
)

// Update сохраняет изменения существующего заказа.
func (r *repository) Update(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.UUID]; !ok {
		return errs.ErrOrderNotFound
	}

	r.orders[order.UUID] = converter.OrderToRecord(order)

	return nil
}
