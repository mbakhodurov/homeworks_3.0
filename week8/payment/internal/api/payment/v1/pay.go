package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/payment/internal/api/converter"
	"github.com/mbakhodurov/homeworks2/week8/payment/internal/service/input"
	payment_v1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	method := converter.PaymentMethodToModel(req.GetPaymentMethod())

	transactionUUID, err := a.paymentService.Pay(ctx, input.PayInput{
		OrderUUID: req.GetOrderUuid(),
		Method:    method,
	})
	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
