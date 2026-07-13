package converter

import (
	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
	orderv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/openapi/order/v1"
)

// ToOrder конвертирует CreateOrderRequest в доменную модель заказа.
// Items на выходе содержит только запрошенные PartUUID — тип и цену
// достраивает сервис через InventoryClient.
func ToOrder(req *orderv1.CreateOrderRequest) model.Order {
	items := []model.OrderItem{
		{PartUUID: req.GetHullUUID()},
		{PartUUID: req.GetEngineUUID()},
	}

	if shieldUUID, ok := req.GetShieldUUID().Get(); ok {
		items = append(items, model.OrderItem{PartUUID: shieldUUID})
	}
	if weaponUUID, ok := req.GetWeaponUUID().Get(); ok {
		items = append(items, model.OrderItem{PartUUID: weaponUUID})
	}

	return model.Order{Items: items}
}

// OrderToDTO конвертирует доменную модель заказа в DTO.
func OrderToDTO(o model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  o.UUID,
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

// OrdersToDTO конвертирует список доменных моделей заказов в DTO.
func OrdersToDTO(orders []model.Order) []orderv1.OrderDto {
	dto := make([]orderv1.OrderDto, 0, len(orders))
	for _, o := range orders {
		dto = append(dto, *OrderToDTO(o))
	}

	return dto
}

// PaymentMethodToModel конвертирует proto-способ оплаты в доменный.
func PaymentMethodToModel(m orderv1.PaymentMethod) model.PaymentMethod {
	return model.PaymentMethod(m)
}
