package converter

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/repository/record"
)

// OrderToRecord конвертирует доменную модель в persistence-shape (без items).
func OrderToRecord(o model.Order) record.Order {
	rec := record.Order{
		UUID:       o.UUID.String(),
		UserUUID:   o.UserUUID.String(),
		TotalPrice: o.TotalPrice(),
		Status:     string(o.Status),
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}

	if o.TransactionUUID != nil {
		s := o.TransactionUUID.String()
		rec.TransactionUUID = &s
	}
	if o.PaymentMethod != nil {
		s := string(*o.PaymentMethod)
		rec.PaymentMethod = &s
	}

	return rec
}

// OrderItemsToRecords конвертирует позиции заказа в persistence-shape.
func OrderItemsToRecords(items []model.OrderItem) []record.OrderItem {
	recs := make([]record.OrderItem, 0, len(items))
	for _, item := range items {
		recs = append(recs, record.OrderItem{
			UUID:      item.UUID.String(),
			OrderUUID: item.OrderUUID.String(),
			PartUUID:  item.PartUUID.String(),
			PartType:  string(item.PartType),
			Price:     item.Price,
		})
	}
	return recs
}

// OrderToModel восстанавливает доменную модель из persistence-shape.
func OrderToModel(r record.Order, items []record.OrderItem) model.Order {
	modelItems := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, model.OrderItem{
			UUID:      uuid.MustParse(item.UUID),
			OrderUUID: uuid.MustParse(item.OrderUUID),
			PartUUID:  uuid.MustParse(item.PartUUID),
			PartType:  model.PartType(item.PartType),
			Price:     item.Price,
		})
	}

	o := model.Order{
		UUID:      uuid.MustParse(r.UUID),
		UserUUID:  uuid.MustParse(r.UserUUID),
		Items:     modelItems,
		Status:    model.OrderStatus(r.Status),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if r.TransactionUUID != nil {
		id := uuid.MustParse(*r.TransactionUUID)
		o.TransactionUUID = &id
	}
	if r.PaymentMethod != nil {
		pm := model.PaymentMethod(*r.PaymentMethod)
		o.PaymentMethod = &pm
	}

	return o
}
