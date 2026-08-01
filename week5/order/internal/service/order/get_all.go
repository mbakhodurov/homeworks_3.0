package order

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week5/order/internal/model"
)

// GetAll возвращает все заказы.
func (s *service) GetAll(ctx context.Context) ([]model.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("получить список заказов: %w", err)
	}

	return orders, nil
}
