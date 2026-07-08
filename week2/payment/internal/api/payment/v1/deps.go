package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week2/payment/internal/service/payment/input"
)

// PaymentService определяет контракт для бизнес-логики оплаты.
type PaymentService interface {
	Pay(ctx context.Context, in input.PayInput) (string, error)
}
