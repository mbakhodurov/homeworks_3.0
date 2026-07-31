package v1

import payment_v1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/proto/payment/v1"

type api struct {
	payment_v1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

// New создаёт gRPC-обработчик PaymentService.
func New(paymentService PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
