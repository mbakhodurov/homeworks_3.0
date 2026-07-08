package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week2/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/repository/converter"
)

// Get возвращает заказ по UUID.
func (r *repository) Get(_ context.Context, orderUUID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.orders[orderUUID]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.OrderToModel(rec), nil
}
