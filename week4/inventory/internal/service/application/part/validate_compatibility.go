package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/input"
)

// ValidateCompatibility проверяет соответствие типов деталей слотам корабля
// и совместимость деталей друг с другом.
func (s *service) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	resolved, err := s.resolveShipSlots(ctx, slots)
	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolved)
}

// resolveShipSlots загружает детали по UUID слотов одним вызовом репозитория,
// проверяет соответствие типа детали слоту и уникальность UUID между слотами.
func (s *service) resolveShipSlots(ctx context.Context, slots model.ShipSlots) (model.ResolvedShipSlots, error) {
	uuids := []uuid.UUID{slots.HullUUID, slots.EngineUUID}
	if slots.ShieldUUID != nil {
		uuids = append(uuids, *slots.ShieldUUID)
	}
	if slots.WeaponUUID != nil {
		uuids = append(uuids, *slots.WeaponUUID)
	}

	if hasDuplicateUUID(uuids) {
		return model.ResolvedShipSlots{}, errs.ErrPartTypeMismatch
	}

	parts, err := s.partRepo.List(ctx, input.PartFilter{UUIDs: uuids})
	if err != nil {
		return model.ResolvedShipSlots{}, fmt.Errorf("получить детали: %w", err)
	}

	byUUID := make(map[uuid.UUID]*model.Part, len(parts))
	for i := range parts {
		byUUID[parts[i].UUID()] = &parts[i]
	}

	hull, err := resolveSlot(byUUID, slots.HullUUID, model.PartTypeHull)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}

	engine, err := resolveSlot(byUUID, slots.EngineUUID, model.PartTypeEngine)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}

	resolved := model.ResolvedShipSlots{Hull: hull, Engine: engine}

	if slots.ShieldUUID != nil {
		shield, shieldErr := resolveSlot(byUUID, *slots.ShieldUUID, model.PartTypeShield)
		if shieldErr != nil {
			return model.ResolvedShipSlots{}, shieldErr
		}
		resolved.Shield = shield
	}

	if slots.WeaponUUID != nil {
		weapon, weaponErr := resolveSlot(byUUID, *slots.WeaponUUID, model.PartTypeWeapon)
		if weaponErr != nil {
			return model.ResolvedShipSlots{}, weaponErr
		}
		resolved.Weapon = weapon
	}

	return resolved, nil
}

// resolveSlot находит деталь по UUID и проверяет, что её тип соответствует ожидаемому слоту.
func resolveSlot(byUUID map[uuid.UUID]*model.Part, partUUID uuid.UUID, expected model.PartType) (*model.Part, error) {
	p, ok := byUUID[partUUID]
	if !ok {
		return nil, errs.ErrPartNotFound
	}

	if p.PartType() != expected {
		return nil, errs.ErrPartTypeMismatch
	}

	return p, nil
}

// hasDuplicateUUID проверяет, что один и тот же UUID не передан сразу в несколько слотов.
func hasDuplicateUUID(uuids []uuid.UUID) bool {
	seen := make(map[uuid.UUID]struct{}, len(uuids))
	for _, u := range uuids {
		if _, ok := seen[u]; ok {
			return true
		}
		seen[u] = struct{}{}
	}

	return false
}
