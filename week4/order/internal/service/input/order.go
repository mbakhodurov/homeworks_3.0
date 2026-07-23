package input

import "github.com/google/uuid"

// CreateOrderInput задаёт параметры создания заказа.
// Космический корабль обязательно состоит из корпуса и двигателя. Щит и оружие —
// опциональные компоненты, поэтому ShieldUUID и WeaponUUID — *uuid.UUID.
type CreateOrderInput struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

// PartUUIDs возвращает UUID всех запрошенных компонентов, включая только non-nil
// опциональные слоты.
func (in CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{in.HullUUID, in.EngineUUID}

	if in.ShieldUUID != nil {
		uuids = append(uuids, *in.ShieldUUID)
	}
	if in.WeaponUUID != nil {
		uuids = append(uuids, *in.WeaponUUID)
	}

	return uuids
}
