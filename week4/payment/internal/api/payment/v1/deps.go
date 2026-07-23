package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week4/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
}
