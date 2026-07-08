package input

import "github.com/mbakhodurov/homeworks2/week2/payment/internal/model"

// PayInput — вход use case'а оплаты заказа
type PayInput struct {
	OrderUUID     string
	PaymentMethod model.PaymentMethod
}
