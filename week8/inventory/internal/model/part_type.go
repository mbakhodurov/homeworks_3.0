package model

import (
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week8/inventory/internal/errors"
)

// PartType — тип детали космического корабля.
type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

// NewPartType создаёт тип детали с валидацией. PartTypeUnspecified не является
// допустимым типом хранимой детали — он используется только как признак «без фильтра».
func NewPartType(s string) (PartType, error) {
	pt := PartType(s)

	switch pt {
	case PartTypeHull, PartTypeEngine, PartTypeShield, PartTypeWeapon:
		return pt, nil
	default:
		return "", fmt.Errorf("неизвестный тип детали %q: %w", s, errs.ErrInvalidProperties)
	}
}
