package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	paymentv1 "github.com/mbakhodurov/homeworks2/week1/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server реализует gRPC сервис оплаты
type server struct {
	mu sync.RWMutex
	paymentv1.UnimplementedPaymentServiceServer
}

// NewServer создаёт новый экземпляр сервера оплаты
func NewServer() *server {
	return &server{}
}

// PayOrder обрабатывает оплату заказа
func (s *server) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid не может быть пустым")
	}

	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "payment_method не может быть UNSPECIFIED")
	}

	// if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "order_uuid имеет неверный формат UUID")
	// }

	transactionUUID := uuid.New().String()

	slog.Info("оплата прошла успешно",
		"order_uuid", req.GetOrderUuid(),
		"transaction_uuid", transactionUUID,
	)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
