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

// Pay проводит оплату заказа.
func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return uuid.Nil, errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return uuid.Nil, errs.ErrOrderCancelled
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method
	order.UpdatedAt = lo.ToPtr(time.Now())

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return uuid.Nil, fmt.Errorf("обновить заказ: %w", err)
	}

	return transactionUUID, nil
}
