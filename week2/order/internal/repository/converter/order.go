package converter

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/repository/record"
)

// OrderToRecord конвертирует доменную модель в persistence-shape.
func OrderToRecord(o model.Order) record.Order {
	items := make([]record.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, record.OrderItem{
			PartUUID: item.PartUUID.String(),
			PartType: string(item.PartType),
			Price:    item.Price,
		})
	}

	rec := record.Order{
		UUID:      o.UUID.String(),
		Items:     items,
		Status:    string(o.Status),
		CreatedAt: o.CreatedAt,
	}

	if o.TransactionUUID != nil {
		transactionUUID := o.TransactionUUID.String()
		rec.TransactionUUID = &transactionUUID
	}
	if o.PaymentMethod != nil {
		paymentMethod := string(*o.PaymentMethod)
		rec.PaymentMethod = &paymentMethod
	}

	return rec
}

// OrderToModel конвертирует запись хранилища в доменную модель.
func OrderToModel(r record.Order) model.Order {
	items := make([]model.OrderItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, model.OrderItem{
			PartUUID: uuid.MustParse(item.PartUUID),
			PartType: model.PartType(item.PartType),
			Price:    item.Price,
		})
	}

	o := model.Order{
		UUID:      uuid.MustParse(r.UUID),
		Items:     items,
		Status:    model.OrderStatus(r.Status),
		CreatedAt: r.CreatedAt,
	}

	if r.TransactionUUID != nil {
		transactionUUID := uuid.MustParse(*r.TransactionUUID)
		o.TransactionUUID = &transactionUUID
	}
	if r.PaymentMethod != nil {
		paymentMethod := model.PaymentMethod(*r.PaymentMethod)
		o.PaymentMethod = &paymentMethod
	}

	return o
}
