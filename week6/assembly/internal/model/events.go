package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderPaidEvent struct {
	EventUUID uuid.UUID
	OrderUUID uuid.UUID
	UserUUID  uuid.UUID
}

type ShipAssembledEvent struct {
	EventUUID    uuid.UUID
	OrderUUID    uuid.UUID
	UserUUID     uuid.UUID
	BuildTimeSec int64
	AssembledAt  time.Time
}
