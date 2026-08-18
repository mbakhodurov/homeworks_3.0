package app

import (
	payapiv1 "github.com/mbakhodurov/homeworks2/week7/payment/internal/api/payment/v1"
	paysvc "github.com/mbakhodurov/homeworks2/week7/payment/internal/service/payment"
	payment_v1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/payment/v1"
)

// diContainer — контейнер зависимостей с ленивой инициализацией.
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует результат.
type diContainer struct {
	// Сервисы
	paymentService payapiv1.PaymentService

	// API-обработчики
	paymentHandler payment_v1.PaymentServiceServer
}

// PaymentService возвращает сервис бизнес-логики оплаты.
func (d *diContainer) PaymentService() payapiv1.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paysvc.New()
	}

	return d.paymentService
}

// PaymentAPI возвращает gRPC-обработчик PaymentService.
func (d *diContainer) PaymentAPI() payment_v1.PaymentServiceServer {
	if d.paymentHandler == nil {
		d.paymentHandler = payapiv1.New(d.PaymentService())
	}

	return d.paymentHandler
}
