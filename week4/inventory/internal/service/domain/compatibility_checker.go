package domain

import (
	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
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
// Соответствие типа детали слоту (правило 1) проверяется раньше, в application-слое,
// на этапе разрешения слотов — сюда приходят уже валидированные по типу детали.
//
// Правила, проверяемые здесь:
//  2. Прочность корпуса должна выдерживать нагрузку двигателя (hull.strength >= engine.requiredStrength).
//  3. Плазменный щит конфликтует с лазерным оружием.
func (c *CompatibilityChecker) Check(slots model.ResolvedShipSlots) error {
	hull := slots.Hull.Properties().Hull()
	engine := slots.Engine.Properties().Engine()

	if !hull.CanSupport(engine) {
		return errs.ErrIncompatibleParts
	}

	if slots.Shield != nil && slots.Weapon != nil {
		shield := slots.Shield.Properties().Shield()
		weapon := slots.Weapon.Properties().Weapon()

		if shield.ConflictsWith(weapon) {
			return errs.ErrIncompatibleParts
		}
	}

	return nil
}
