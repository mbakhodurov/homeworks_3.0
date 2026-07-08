package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/mbakhodurov/homeworks2/week2/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
)

// Cancel отменяет заказ, который ещё не был оплачен.
func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderCancelled
	}

	order.Status = model.OrderStatusCancelled
	order.UpdatedAt = lo.ToPtr(time.Now())

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	return nil
}
