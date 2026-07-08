package order

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/repository/converter"
)

// GetAll возвращает все заказы без фильтрации.
func (r *repository) GetAll(_ context.Context) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]model.Order, 0, len(r.orders))
	for _, rec := range r.orders {
		orders = append(orders, converter.OrderToModel(rec))
	}

	return orders, nil
}
