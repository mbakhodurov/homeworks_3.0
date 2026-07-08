package record

import "time"

// Order — persistence-shape заказа для хранилища.
type Order struct {
	UUID            string
	Items           []OrderItem
	TransactionUUID *string
	PaymentMethod   *string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

// OrderItem — persistence-shape позиции заказа.
type OrderItem struct {
	PartUUID string
	PartType string
	Price    int64
}
