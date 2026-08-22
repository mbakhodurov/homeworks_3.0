package converter

import (
	"github.com/mbakhodurov/homeworks2/week8/order/internal/model"
	paymentv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/payment/v1"
)

// PaymentMethodToDTO конвертирует доменный способ оплаты в proto-представление.
func PaymentMethodToDTO(m model.PaymentMethod) paymentv1.PaymentMethod {
	switch m {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
