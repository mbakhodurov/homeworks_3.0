package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/mbakhodurov/homeworks2/week4/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/order/internal/model"
)

// Pay проводит оплату заказа.
//
// gRPC-вызов к PaymentService выполняется внутри открытой транзакции — учебное упрощение.
// В продакшене это плохо: если PaymentService отвечает медленно, транзакция держится долго,
// а если сервис упал уже после списания денег, транзакция откатится, хотя оплата прошла.
// Проблемы распределённых транзакций разберём в следующих неделях.
func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	var transactionUUID uuid.UUID

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.Get(txCtx, orderUUID)
		if err != nil {
			return fmt.Errorf("получить заказ: %w", err)
		}

		if order.Status != model.OrderStatusPendingPayment {
			if order.Status == model.OrderStatusPaid {
				return errs.ErrOrderAlreadyPaid
			}
			return errs.ErrOrderCancelled
		}

		transactionUUID, err = s.paymentClient.PayOrder(txCtx, orderUUID, method)
		if err != nil {
			return fmt.Errorf("оплатить заказ: %w", err)
		}

		order.Status = model.OrderStatusPaid
		order.TransactionUUID = &transactionUUID
		order.PaymentMethod = &method
		order.UpdatedAt = lo.ToPtr(time.Now())

		if err := s.orderRepo.Update(txCtx, order); err != nil {
			return fmt.Errorf("обновить заказ: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
