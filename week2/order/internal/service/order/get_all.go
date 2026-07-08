package order

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
)

// GetAll возвращает все заказы без фильтрации.
func (s *service) GetAll(ctx context.Context) ([]model.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("получить все заказы: %w", err)
	}

	return orders, nil
}
