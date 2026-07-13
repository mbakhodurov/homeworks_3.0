package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/client/grpc/payment/v1/converter"
	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
	paymentv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/payment/v1"
)

// client — обёртка над proto-клиентом PaymentService.
type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

// New создаёт новый клиент PaymentService.
func New(paymentClient paymentv1.PaymentServiceClient) *client {
	return &client{paymentClient: paymentClient}
}

// PayOrder проводит оплату заказа и возвращает UUID транзакции.
func (c *client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      orderUUID.String(),
		PaymentMethod: converter.PaymentMethodToDTO(method),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("распарсить uuid транзакции: %w", err)
	}

	return transactionUUID, nil
}
