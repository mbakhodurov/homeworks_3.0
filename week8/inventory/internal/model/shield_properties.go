package model

import (
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week8/inventory/internal/errors"
)

// ShieldType — тип щита.
type ShieldType string

const (
	ShieldTypeEnergy ShieldType = "energy"
	ShieldTypePlasma ShieldType = "plasma"
)

// ShieldProperties — свойства щита (Value Object).
type ShieldProperties struct {
	shieldType ShieldType
}

// NewShieldProperties создаёт свойства щита с валидацией типа.
func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	switch shieldType {
	case ShieldTypeEnergy, ShieldTypePlasma:
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип щита %q: %w", shieldType, errs.ErrInvalidProperties)
	}

	return PartProperties{shield: &ShieldProperties{shieldType: shieldType}}, nil
}

// ShieldType возвращает тип щита.
func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

// ConflictsWith проверяет конфликт плазменного щита с лазерным оружием —
// плазма создаёт электромагнитные помехи, нарушающие работу лазера.
func (s *ShieldProperties) ConflictsWith(w *WeaponProperties) bool {
	return s.shieldType == ShieldTypePlasma && w.weaponType == WeaponTypeLaser
}
