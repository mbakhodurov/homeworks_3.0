package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/mbakhodurov/homeworks2/week5/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week5/order/internal/model"
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
		order, err := s.orderRepo.GetForUpdate(txCtx, orderUUID)
		if err != nil {
			return fmt.Errorf("получить заказ: %w", err)
		}

		switch order.Status {
		case model.OrderStatusPaid:
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			return errs.ErrOrderCancelled
		case model.OrderStatusAssembled:
			return errs.ErrOrderAlreadyAssembled
		case model.OrderStatusPendingPayment:
			// ok
		default:
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

		if err := s.producer.Produce(txCtx, model.OrderPaidEvent{
			EventUUID: uuid.New(),
			OrderUUID: orderUUID,
			UserUUID:  order.UserUUID,
		}); err != nil {
			return fmt.Errorf("отправить событие оплаты: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
