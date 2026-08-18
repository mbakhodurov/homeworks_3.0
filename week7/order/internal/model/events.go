package model

import "github.com/google/uuid"

// OrderPaidEvent — событие оплаты заказа (публикуется OrderService в Kafka).
type OrderPaidEvent struct {
	EventUUID uuid.UUID
	OrderUUID uuid.UUID
	UserUUID  uuid.UUID
}

// ShipAssembledEvent — событие завершения сборки корабля (потребляется OrderService из Kafka).
type ShipAssembledEvent struct {
	EventUUID    uuid.UUID
	OrderUUID    uuid.UUID
	UserUUID     uuid.UUID
	BuildTimeSec int64
}
