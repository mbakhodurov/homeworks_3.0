package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week3/payment/internal/model"
)

// PaymentService определяет контракт для бизнес-логики оплаты.
type PaymentService interface {
	Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
}
