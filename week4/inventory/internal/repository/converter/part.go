package converter

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/repository/record"
)

// PartRecordToModel конвертирует запись хранилища в доменную модель.
func PartRecordToModel(rec record.Part) (model.Part, error) {
	var propsRec record.PartPropertiesRecord
	if err := json.Unmarshal(rec.Properties, &propsRec); err != nil {
		return model.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
	}

	props, err := partPropertiesFromRecord(propsRec)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать свойства: %w", err)
	}

	partType, err := model.NewPartType(rec.PartType)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	partUUID, err := uuid.Parse(rec.UUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать uuid детали: %w", err)
	}

	return model.RestorePart(
		partUUID,
		rec.Name,
		rec.Description,
		partType,
		rec.Price,
		rec.StockQuantity,
		rec.Reserved,
		props,
		rec.CreatedAt,
	), nil
}

// partPropertiesFromRecord определяет тип свойств по тому, какое поле record non-nil,
// и создаёт соответствующий Value Object.
func partPropertiesFromRecord(rec record.PartPropertiesRecord) (model.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return model.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return model.NewEngineProperties(model.EngineClass(rec.Engine.Class), rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return model.NewShieldProperties(model.ShieldType(rec.Shield.ShieldType))
	case rec.Weapon != nil:
		return model.NewWeaponProperties(model.WeaponType(rec.Weapon.WeaponType))
	default:
		return model.PartProperties{}, nil
	}
}
