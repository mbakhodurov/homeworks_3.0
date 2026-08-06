package model

import "github.com/google/uuid"

// Part — то, что order получил из inventory по gRPC. Отдельная от
// inventory/internal/model.Part сущность: order знает только то, что нужно
// ему самому, и ничего не знает про proto-типы inventory.
type Part struct {
	UUID          uuid.UUID
	Name          string
	PartType      PartType
	Price         int64
	StockQuantity int64
}

type PartType string

const (
	PartTypeHull   PartType = "HULL"
	PartTypeEngine PartType = "ENGINE"
	PartTypeShield PartType = "SHIELD"
	PartTypeWeapon PartType = "WEAPON"
)

// ShipSlots — UUID-ы деталей по слотам корабля. Hull и Engine — обязательные слоты,
// Shield и Weapon — опциональные (nil означает, что слот не используется).
type ShipSlots struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}
