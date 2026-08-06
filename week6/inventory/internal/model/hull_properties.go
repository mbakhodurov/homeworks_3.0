package model

import (
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week6/inventory/internal/errors"
)

const (
	minHullStrength = 30
	maxHullStrength = 200
)

// HullProperties — свойства корпуса (Value Object).
type HullProperties struct {
	strength int
}

// NewHullProperties создаёт свойства корпуса. Прочность должна быть в диапазоне 30–200.
func NewHullProperties(strength int) (PartProperties, error) {
	if strength < minHullStrength || strength > maxHullStrength {
		return PartProperties{}, fmt.Errorf("прочность корпуса должна быть от 30 до 200, получено %d: %w", strength, errs.ErrInvalidProperties)
	}

	return PartProperties{hull: &HullProperties{strength: strength}}, nil
}

// Strength возвращает прочность корпуса.
func (h *HullProperties) Strength() int { return h.strength }

// CanSupport проверяет, выдержит ли корпус нагрузку двигателя.
func (h *HullProperties) CanSupport(e *EngineProperties) bool {
	return h.strength >= e.requiredStrength
}
