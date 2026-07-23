package app

import (
	payapiv1 "github.com/mbakhodurov/homeworks2/week4/payment/internal/api/payment/v1"
	"github.com/mbakhodurov/homeworks2/week4/payment/internal/interceptor"
	paysvc "github.com/mbakhodurov/homeworks2/week4/payment/internal/service/payment"
	payment_v1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

// Interceptors возвращает gRPC-опции с интерцепторами для использования в тестах.
func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor(),
			interceptor.ErrorInterceptor,
			interceptor.LoggerInterceptor(),
		),
	}
}

// RegisterServices собирает зависимости и регистрирует gRPC-сервисы на сервере.
func RegisterServices(grpcServer *grpc.Server) {
	svc := paysvc.New()
	api := payapiv1.New(svc)
	payment_v1.RegisterPaymentServiceServer(grpcServer, api)
}
