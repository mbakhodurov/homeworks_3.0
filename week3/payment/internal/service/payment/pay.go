package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	errs "github.com/mbakhodurov/homeworks2/week3/payment/internal/errors"
	"github.com/mbakhodurov/homeworks2/week3/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	if orderUUID == "" {
		return "", errs.ErrInvalidOrderUUID
	}
	if _, err := uuid.Parse(orderUUID); err != nil {
		return "", errs.ErrInvalidOrderUUID
	}
	if !method.IsValid() {
		return "", errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New().String()

	slog.InfoContext(ctx, "оплата выполнена",
		"order_uuid", orderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", method)

	return transactionUUID, nil
}
