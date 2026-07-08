package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks2/week2/payment/internal/api/converter"
	errs "github.com/mbakhodurov/homeworks2/week2/payment/internal/errors"
	payment_v1 "github.com/mbakhodurov/homeworks2/week2/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	method := converter.PaymentMethodToModel(req.GetPaymentMethod())

	transactionUUID, err := a.paymentService.Pay(ctx, req.GetOrderUuid(), method)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidOrderUUID) || errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
