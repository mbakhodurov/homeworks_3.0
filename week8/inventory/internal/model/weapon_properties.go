package model

import (
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week8/inventory/internal/errors"
)

// WeaponType — тип оружия.
type WeaponType string

const (
	WeaponTypeLaser   WeaponType = "laser"
	WeaponTypeMissile WeaponType = "missile"
)

// WeaponProperties — свойства оружия (Value Object).
type WeaponProperties struct {
	weaponType WeaponType
}

// NewWeaponProperties создаёт свойства оружия с валидацией типа.
func NewWeaponProperties(weaponType WeaponType) (PartProperties, error) {
	switch weaponType {
	case WeaponTypeLaser, WeaponTypeMissile:
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип оружия %q: %w", weaponType, errs.ErrInvalidProperties)
	}

	return PartProperties{weapon: &WeaponProperties{weaponType: weaponType}}, nil
}

// WeaponType возвращает тип оружия.
func (w *WeaponProperties) WeaponType() WeaponType { return w.weaponType }
