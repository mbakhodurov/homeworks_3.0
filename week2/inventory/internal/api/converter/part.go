package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/mbakhodurov/homeworks2/week2/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week2/shared/pkg/proto/inventory/v1"
)

// ToUUID парсит строковый UUID из транспортного запроса.
func ToUUID(rawUUID string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return parsed, nil
}

// ToListFilter парсит фильтр ListPartsRequest в uuids и тип детали.
//
// PartsFilter.part_type в proto — список (repeated), а доменная фильтрация
// поддерживает один тип: берём первый элемент списка, если он есть.
func ToListFilter(req *inventoryv1.ListPartsRequest) ([]uuid.UUID, model.PartType, error) {
	filter := req.GetFilter()
	if filter == nil {
		return nil, model.PartTypeUnspecified, nil
	}

	uuids := make([]uuid.UUID, 0, len(filter.GetUuids()))
	for _, rawUUID := range filter.GetUuids() {
		parsed, err := ToUUID(rawUUID)
		if err != nil {
			return nil, model.PartTypeUnspecified, err
		}
		uuids = append(uuids, parsed)
	}

	var partType model.PartType
	if types := filter.GetPartType(); len(types) > 0 {
		partType = PartTypeToModel(types[0])
	}

	return uuids, partType, nil
}

// ToPart конвертирует CreatePartsRequest в доменную модель детали.
// UUID и CreatedAt не заполняются — их назначает сервис при создании.
func ToPart(req *inventoryv1.CreatePartsRequest) (model.Part, error) {
	info := req.GetInfo()
	if info == nil {
		return model.Part{}, errs.ErrInvalidPartInfo
	}

	return model.Part{
		Name:          info.GetName(),
		Description:   info.GetDescription(),
		Price:         info.GetPrice(),
		PartType:      PartTypeToModel(info.GetPartType()),
		StockQuantity: info.GetStockQuantity(),
	}, nil
}

// PartToDTO конвертирует доменную модель детали в proto-представление.
func PartToDTO(p model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid: p.UUID.String(),
		Info: &inventoryv1.PartInfo{
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			PartType:      PartTypeToDTO(p.PartType),
			StockQuantity: p.StockQuantity,
		},
		CreatedAt: timestamppb.New(p.CreatedAt),
	}
}

// PartsToDTO конвертирует список доменных моделей в proto-представление.
func PartsToDTO(parts []model.Part) []*inventoryv1.Part {
	dto := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		dto = append(dto, PartToDTO(p))
	}

	return dto
}

// PartTypeToModel конвертирует proto-тип детали в доменный.
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
		return model.PartTypeUnspecified
	}
}

// PartTypeToDTO конвертирует доменный тип детали в proto-представление.
func PartTypeToDTO(t model.PartType) inventoryv1.PartType {
	switch t {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}
