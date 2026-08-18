package converter

import (
	"github.com/mbakhodurov/homeworks2/week7/payment/internal/model"
	payment_v1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/payment/v1"
)

func PaymentMethodToModel(method payment_v1.PaymentMethod) model.PaymentMethod {
	switch method {
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case payment_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
