package model

import (
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week6/inventory/internal/errors"
)

// EngineClass — класс двигателя.
type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

// EngineProperties — свойства двигателя (Value Object).
type EngineProperties struct {
	class            EngineClass
	requiredStrength int
}

// NewEngineProperties создаёт свойства двигателя. Класс должен быть A/B/C,
// требуемая прочность корпуса — положительной.
func NewEngineProperties(class EngineClass, requiredStrength int) (PartProperties, error) {
	switch class {
	case EngineClassA, EngineClassB, EngineClassC:
	default:
		return PartProperties{}, fmt.Errorf("неизвестный класс двигателя %q: %w", class, errs.ErrInvalidProperties)
	}

	if requiredStrength <= 0 {
		return PartProperties{}, fmt.Errorf(
			"требуемая прочность корпуса должна быть положительной: %w", errs.ErrInvalidProperties,
		)
	}

	return PartProperties{engine: &EngineProperties{class: class, requiredStrength: requiredStrength}}, nil
}

// Class возвращает класс двигателя.
func (e *EngineProperties) Class() EngineClass { return e.class }

// RequiredStrength возвращает минимальную прочность корпуса, необходимую для этого двигателя.
func (e *EngineProperties) RequiredStrength() int { return e.requiredStrength }
