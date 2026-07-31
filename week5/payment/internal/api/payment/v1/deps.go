package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week5/payment/internal/service/input"
)

// PaymentService — бизнес-логика оплаты, которую использует gRPC-обработчик.
type PaymentService interface {
	Pay(ctx context.Context, in input.PayInput) (string, error)
}
