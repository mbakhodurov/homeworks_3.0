package input

import "github.com/mbakhodurov/homeworks2/week5/payment/internal/model"

// PayInput — вход use case'а оплаты заказа.
type PayInput struct {
	OrderUUID string
	Method    model.PaymentMethod
}
