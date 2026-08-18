package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week7/payment/internal/errors"
	"github.com/mbakhodurov/homeworks2/week7/payment/internal/service/input"
)

func (s *service) Pay(ctx context.Context, in input.PayInput) (string, error) {
	if in.OrderUUID == "" {
		return "", errs.ErrInvalidOrderUUID
	}
	if _, err := uuid.Parse(in.OrderUUID); err != nil {
		return "", errs.ErrInvalidOrderUUID
	}
	if !in.Method.IsValid() {
		return "", errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New().String()

	slog.InfoContext(ctx, "оплата выполнена",
		"order_uuid", in.OrderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", in.Method)

	return transactionUUID, nil
}
