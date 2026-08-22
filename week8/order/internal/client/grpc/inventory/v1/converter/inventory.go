package converter

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week8/order/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/inventory/v1"
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

// ToUUIDStrings конвертирует список UUID в список строк для proto-запросов
// ReserveParts/ReleaseParts.
func ToUUIDStrings(uuids []uuid.UUID) []string {
	rawUUIDs := make([]string, 0, len(uuids))
	for _, id := range uuids {
		rawUUIDs = append(rawUUIDs, id.String())
	}

	return rawUUIDs
}

// ShipSlotsToProto конвертирует доменные слоты корабля в proto-запрос ValidateCompatibility.
func ShipSlotsToProto(slots model.ShipSlots) *inventoryv1.ValidateCompatibilityRequest {
	req := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   slots.HullUUID.String(),
		EngineUuid: slots.EngineUUID.String(),
	}

	if slots.ShieldUUID != nil {
		req.ShieldUuid = slots.ShieldUUID.String()
	}
	if slots.WeaponUUID != nil {
		req.WeaponUuid = slots.WeaponUUID.String()
	}

	return req
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
