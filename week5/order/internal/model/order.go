package model

import (
	"time"

	"github.com/google/uuid"
)

// Order — aggregate root заказа на постройку космического корабля.
// Остаётся анемичной моделью (публичные поля) — DDD для заказов будет в следующих неделях.
// Содержит позиции Items, каждая из которых ссылается на деталь через PartUUID.
type Order struct {
	UUID            uuid.UUID
	Items           []OrderItem
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

// TotalPrice возвращает сумму цен всех позиций заказа.
func (o Order) TotalPrice() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
	OrderStatusAssembled      OrderStatus = "ASSEMBLED"
)

type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// CreateOrderInput задаёт параметры создания заказа.
// Hull и Engine — обязательные слоты, Shield и Weapon — опциональные.
type CreateOrderInput struct {
	UserUUID   uuid.UUID
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

// PartUUIDs возвращает UUID всех запрошенных компонентов, включая только non-nil опциональные слоты.
func (in CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{in.HullUUID, in.EngineUUID}

	if in.ShieldUUID != nil {
		uuids = append(uuids, *in.ShieldUUID)
	}
	if in.WeaponUUID != nil {
		uuids = append(uuids, *in.WeaponUUID)
	}

	return uuids
}
