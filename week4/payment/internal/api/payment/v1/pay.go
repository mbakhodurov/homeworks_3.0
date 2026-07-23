package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week4/payment/internal/converter"
	payment_v1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	method := converter.PaymentMethodToModel(req.GetPaymentMethod())

	transactionUUID, err := a.paymentService.Pay(ctx, req.GetOrderUuid(), method)
	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
