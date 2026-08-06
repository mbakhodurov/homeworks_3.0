package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
)

// Get возвращает заказ по UUID.
func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return order, nil
}
