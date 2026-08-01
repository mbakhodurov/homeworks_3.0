package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week5/payment/internal/service/input"
)

type PaymentService interface {
	Pay(ctx context.Context, in input.PayInput) (string, error)
}
