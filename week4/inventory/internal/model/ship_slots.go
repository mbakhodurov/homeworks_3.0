package model

import "github.com/google/uuid"

// ShipSlots — UUID-ы деталей по слотам корабля, до проверки типов.
// Hull и Engine — обязательные слоты, Shield и Weapon — опциональные
// (nil означает, что слот не используется).
type ShipSlots struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

// ResolvedShipSlots — слоты корабля с уже разрешёнными и проверенными по типу деталями.
// Hull и Engine гарантированно заполнены, Shield и Weapon могут быть nil.
type ResolvedShipSlots struct {
	Hull   *Part
	Engine *Part
	Shield *Part
	Weapon *Part
}
