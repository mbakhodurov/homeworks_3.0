package input

import "github.com/mbakhodurov/homeworks2/week8/payment/internal/model"

type PayInput struct {
	OrderUUID string
	Method    model.PaymentMethod
}
