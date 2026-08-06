package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/mbakhodurov/homeworks2/week6/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
)

// Cancel отменяет заказ, который ещё не был оплачен.
//
// Без транзакции — осознанное упрощение: сначала освобождаем резерв (ReleaseParts, gRPC),
// затем обновляем статус заказа. Если Update упадёт после успешного Release, резерв будет
// снят, а заказ останется в PENDING_PAYMENT. Корректное решение — saga/outbox.
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

	uuids := make([]uuid.UUID, 0, len(order.Items))
	for _, item := range order.Items {
		uuids = append(uuids, item.PartUUID)
	}

	if err := s.inventoryClient.ReleaseParts(ctx, uuids); err != nil {
		return fmt.Errorf("освободить детали: %w", err)
	}

	order.Status = model.OrderStatusCancelled
	order.UpdatedAt = lo.ToPtr(time.Now())

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	return nil
}
