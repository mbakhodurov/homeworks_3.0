package converter

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/inventory/v1"
)

// ToListPartsRequest конвертирует список UUID в proto-запрос ListParts.
func ToListPartsRequest(uuids []uuid.UUID) *inventoryv1.ListPartsRequest {
	rawUUIDs := make([]string, 0, len(uuids))
	for _, id := range uuids {
		rawUUIDs = append(rawUUIDs, id.String())
	}

	return &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{Uuids: rawUUIDs},
	}
}

// PartsToModel конвертирует proto-детали в доменную модель order.
func PartsToModel(parts []*inventoryv1.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, PartToModel(p))
	}

	return result
}

// PartToModel конвертирует одну proto-деталь в доменную модель order.
func PartToModel(p *inventoryv1.Part) model.Part {
	return model.Part{
		UUID:          uuid.MustParse(p.GetUuid()),
		Name:          p.GetInfo().GetName(),
		PartType:      PartTypeToModel(p.GetInfo().GetPartType()),
		Price:         p.GetInfo().GetPrice(),
		StockQuantity: p.GetInfo().GetStockQuantity(),
	}
}

// PartTypeToModel конвертирует proto-тип детали в доменный тип order.
func PartTypeToModel(t inventoryv1.PartType) model.PartType {
	switch t {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return ""
	}
}
