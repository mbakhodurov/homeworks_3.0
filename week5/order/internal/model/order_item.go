package model

import "github.com/google/uuid"

// OrderItem — позиция заказа (снимок детали на момент создания).
type OrderItem struct {
	UUID      uuid.UUID
	OrderUUID uuid.UUID
	PartUUID  uuid.UUID
	PartType  PartType
	Price     int64
}
