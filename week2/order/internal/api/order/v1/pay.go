package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/api/converter"
	orderv1 "github.com/mbakhodurov/homeworks2/week2/shared/pkg/openapi/order/v1"
)

// PayOrder проводит оплату заказа.
func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	method := converter.PaymentMethodToModel(req.GetPaymentMethod())

	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, method)
	if err != nil {
		return nil, err
	}

	return &orderv1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}
