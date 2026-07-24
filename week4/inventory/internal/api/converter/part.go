package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/input"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// ToUUID парсит строковый UUID из транспортного запроса.
func ToUUID(rawUUID string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return parsed, nil
}

// ToUUIDs парсит список строковых UUID.
func ToUUIDs(rawUUIDs []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(rawUUIDs))
	for _, raw := range rawUUIDs {
		parsed, err := ToUUID(raw)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, parsed)
	}

	return uuids, nil
}

// ToListFilter парсит фильтр ListPartsRequest в PartFilter.
//
// PartsFilter.part_type в proto — список (repeated), а доменная фильтрация
// поддерживает один тип: берём первый элемент списка, если он есть.
func ToListFilter(req *inventoryv1.ListPartsRequest) (input.PartFilter, error) {
	filter := req.GetFilter()
	if filter == nil {
		return input.PartFilter{}, nil
	}

	uuids, err := ToUUIDs(filter.GetUuids())
	if err != nil {
		return input.PartFilter{}, err
	}

	var partType model.PartType
	if types := filter.GetPartType(); len(types) > 0 {
		partType = PartTypeToModel(types[0])
	}

	return input.PartFilter{UUIDs: uuids, PartType: partType}, nil
}

// ToCreatePartInput конвертирует CreatePartsRequest во входные данные создания детали.
func ToCreatePartInput(req *inventoryv1.CreatePartsRequest) (input.CreatePartInput, error) {
	info := req.GetInfo()
	if info == nil {
		return input.CreatePartInput{}, errs.ErrInvalidPartInfo
	}

	return input.CreatePartInput{
		Name:          info.GetName(),
		Description:   info.GetDescription(),
		Price:         info.GetPrice(),
		PartType:      PartTypeToModel(info.GetPartType()),
		StockQuantity: int(info.GetStockQuantity()),
	}, nil
}

// ToShipSlots конвертирует ValidateCompatibilityRequest в доменные слоты корабля.
func ToShipSlots(req *inventoryv1.ValidateCompatibilityRequest) (model.ShipSlots, error) {
	hullUUID, err := ToUUID(req.GetHullUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	engineUUID, err := ToUUID(req.GetEngineUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	slots := model.ShipSlots{HullUUID: hullUUID, EngineUUID: engineUUID}

	if raw := req.GetShieldUuid(); raw != "" {
		shieldUUID, shieldErr := ToUUID(raw)
		if shieldErr != nil {
			return model.ShipSlots{}, shieldErr
		}
		slots.ShieldUUID = &shieldUUID
	}

	if raw := req.GetWeaponUuid(); raw != "" {
		weaponUUID, weaponErr := ToUUID(raw)
		if weaponErr != nil {
			return model.ShipSlots{}, weaponErr
		}
		slots.WeaponUUID = &weaponUUID
	}

	return slots, nil
}

// PartToDTO конвертирует доменную модель детали в proto-представление.
func PartToDTO(p model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid: p.UUID().String(),
		Info: &inventoryv1.PartInfo{
			Name:          p.Name(),
			Description:   p.Description(),
			Price:         p.Price(),
			PartType:      PartTypeToDTO(p.PartType()),
			StockQuantity: int64(p.StockQuantity()),
		},
		CreatedAt:  timestamppb.New(p.CreatedAt()),
		Properties: PartPropertiesToDTO(p.Properties()),
	}
}

// PartPropertiesToDTO конвертирует типоспецифичные свойства детали в proto-представление.
// Возвращает nil, если ни одно из полей не заполнено (деталь создана без properties).
func PartPropertiesToDTO(props *model.PartProperties) *inventoryv1.PartProperties {
	switch {
	case props.Hull() != nil:
		return &inventoryv1.PartProperties{
			Properties: &inventoryv1.PartProperties_Hull{
				Hull: &inventoryv1.HullProperties{Strength: int64(props.Hull().Strength())},
			},
		}
	case props.Engine() != nil:
		return &inventoryv1.PartProperties{
			Properties: &inventoryv1.PartProperties_Engine{
				Engine: &inventoryv1.EngineProperties{
					Class:            EngineClassToDTO(props.Engine().Class()),
					RequiredStrength: int64(props.Engine().RequiredStrength()),
				},
			},
		}
	case props.Shield() != nil:
		return &inventoryv1.PartProperties{
			Properties: &inventoryv1.PartProperties_Shield{
				Shield: &inventoryv1.ShieldProperties{ShieldType: ShieldTypeToDTO(props.Shield().ShieldType())},
			},
		}
	case props.Weapon() != nil:
		return &inventoryv1.PartProperties{
			Properties: &inventoryv1.PartProperties_Weapon{
				Weapon: &inventoryv1.WeaponProperties{WeaponType: WeaponTypeToDTO(props.Weapon().WeaponType())},
			},
		}
	default:
		return nil
	}
}

// EngineClassToDTO конвертирует доменный класс двигателя в proto-представление.
func EngineClassToDTO(c model.EngineClass) inventoryv1.EngineClass {
	switch c {
	case model.EngineClassA:
		return inventoryv1.EngineClass_ENGINE_CLASS_A
	case model.EngineClassB:
		return inventoryv1.EngineClass_ENGINE_CLASS_B
	case model.EngineClassC:
		return inventoryv1.EngineClass_ENGINE_CLASS_C
	default:
		return inventoryv1.EngineClass_ENGINE_CLASS_UNSPECIFIED
	}
}

// ShieldTypeToDTO конвертирует доменный тип щита в proto-представление.
func ShieldTypeToDTO(t model.ShieldType) inventoryv1.ShieldType {
	switch t {
	case model.ShieldTypeEnergy:
		return inventoryv1.ShieldType_SHIELD_TYPE_ENERGY
	case model.ShieldTypePlasma:
		return inventoryv1.ShieldType_SHIELD_TYPE_PLASMA
	default:
		return inventoryv1.ShieldType_SHIELD_TYPE_UNSPECIFIED
	}
}

// WeaponTypeToDTO конвертирует доменный тип оружия в proto-представление.
func WeaponTypeToDTO(t model.WeaponType) inventoryv1.WeaponType {
	switch t {
	case model.WeaponTypeLaser:
		return inventoryv1.WeaponType_WEAPON_TYPE_LASER
	case model.WeaponTypeMissile:
		return inventoryv1.WeaponType_WEAPON_TYPE_MISSILE
	default:
		return inventoryv1.WeaponType_WEAPON_TYPE_UNSPECIFIED
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
