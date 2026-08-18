package domain

import (
	errs "github.com/mbakhodurov/homeworks2/week7/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week7/inventory/internal/model"
)

// CompatibilityChecker — доменный сервис проверки совместимости деталей корабля.
// Stateless: не хранит состояние, работает только с переданными данными.
type CompatibilityChecker struct{}

// NewCompatibilityChecker создаёт новый доменный сервис проверки совместимости.
func NewCompatibilityChecker() *CompatibilityChecker {
	return &CompatibilityChecker{}
}

// Check проверяет совместимость деталей корабля.
//
// Правила:
//  1. Прочность корпуса должна выдерживать нагрузку двигателя (hull.strength >= engine.requiredStrength).
//  2. Плазменный щит конфликтует с лазерным оружием.
func (c *CompatibilityChecker) Check(parts []model.Part) error {
	var hull, engine, shield, weapon *model.Part

	for i := range parts {
		switch parts[i].PartType() {
		case model.PartTypeHull:
			hull = &parts[i]
		case model.PartTypeEngine:
			engine = &parts[i]
		case model.PartTypeShield:
			shield = &parts[i]
		case model.PartTypeWeapon:
			weapon = &parts[i]
		}
	}

	if hull == nil || engine == nil {
		return errs.ErrPartNotFound
	}

	if !hull.Properties().Hull().CanSupport(engine.Properties().Engine()) {
		return errs.ErrIncompatibleParts
	}

	if shield != nil && weapon != nil {
		if shield.Properties().Shield().ConflictsWith(weapon.Properties().Weapon()) {
			return errs.ErrIncompatibleParts
		}
	}

	return nil
}
