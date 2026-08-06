package converter

import (
	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
	orderv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/openapi/order/v1"
)

// ToCreateOrderInput конвертирует CreateOrderRequest во входные данные создания заказа.
func ToCreateOrderInput(req *orderv1.CreateOrderRequest) model.CreateOrderInput {
	in := model.CreateOrderInput{
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
	}

	if shieldUUID, ok := req.GetShieldUUID().Get(); ok {
		in.ShieldUUID = &shieldUUID
	}
	if weaponUUID, ok := req.GetWeaponUUID().Get(); ok {
		in.WeaponUUID = &weaponUUID
	}

	return in
}

// OrderToDTO конвертирует доменную модель заказа в DTO.
func OrderToDTO(o model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  o.UUID,
		UserUUID:   o.UserUUID,
		TotalPrice: o.TotalPrice(),
		Status:     orderv1.OrderStatus(o.Status),
		CreatedAt:  o.CreatedAt,
	}

	for _, item := range o.Items {
		switch item.PartType {
		case model.PartTypeHull:
			dto.HullUUID = item.PartUUID
		case model.PartTypeEngine:
			dto.EngineUUID = item.PartUUID
		case model.PartTypeShield:
			dto.ShieldUUID = orderv1.NewOptNilUUID(item.PartUUID)
		case model.PartTypeWeapon:
			dto.WeaponUUID = orderv1.NewOptNilUUID(item.PartUUID)
		}
	}

	if o.TransactionUUID != nil {
		dto.TransactionUUID = orderv1.NewOptNilUUID(*o.TransactionUUID)
	}
	if o.PaymentMethod != nil {
		dto.PaymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*o.PaymentMethod))
	}

	return dto
}

// PaymentMethodToModel конвертирует proto-способ оплаты в доменный.
func PaymentMethodToModel(m orderv1.PaymentMethod) model.PaymentMethod {
	return model.PaymentMethod(m)
}
