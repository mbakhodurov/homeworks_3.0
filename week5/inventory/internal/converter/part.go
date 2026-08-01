package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/mbakhodurov/homeworks2/week5/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week5/inventory/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/proto/inventory/v1"
)

// ToListFilter парсит фильтр ListPartsRequest в PartFilter.
//
// PartsFilter.part_type в proto — список (repeated), а доменная фильтрация
// поддерживает один тип: берём первый элемент списка, если он есть.
func ToListFilter(req *inventoryv1.ListPartsRequest) (model.PartFilter, error) {
	filter := req.GetFilter()
	if filter == nil {
		return model.PartFilter{}, nil
	}

	var partType model.PartType
	if types := filter.GetPartType(); len(types) > 0 {
		partType = PartTypeToModel(types[0])
	}

	return model.PartFilter{UUIDs: filter.GetUuids(), PartType: partType}, nil
}

// ToCreatePartInput конвертирует CreatePartsRequest во входные данные создания детали.
func ToCreatePartInput(req *inventoryv1.CreatePartsRequest) (model.CreatePartInput, error) {
	info := req.GetInfo()
	if info == nil {
		return model.CreatePartInput{}, errs.ErrInvalidPartInfo
	}

	props, err := ToModelProperties(info.GetProperties())
	if err != nil {
		return model.CreatePartInput{}, err
	}

	return model.CreatePartInput{
		Name:          info.GetName(),
		Description:   info.GetDescription(),
		Price:         info.GetPrice(),
		PartType:      PartTypeToModel(info.GetPartType()),
		StockQuantity: int(info.GetStockQuantity()),
		Properties:    props,
	}, nil
}

// ToModelProperties конвертирует proto-свойства детали в доменную модель.
func ToModelProperties(proto *inventoryv1.PartProperties) (model.PartProperties, error) {
	if proto == nil {
		return model.PartProperties{}, nil
	}
	switch p := proto.GetProperties().(type) {
	case *inventoryv1.PartProperties_Hull:
		return model.NewHullProperties(int(p.Hull.GetStrength()))
	case *inventoryv1.PartProperties_Engine:
		return model.NewEngineProperties(
			EngineClassToModel(p.Engine.GetClass()),
			int(p.Engine.GetRequiredStrength()),
		)
	case *inventoryv1.PartProperties_Shield:
		return model.NewShieldProperties(ShieldTypeToModel(p.Shield.GetShieldType()))
	case *inventoryv1.PartProperties_Weapon:
		return model.NewWeaponProperties(WeaponTypeToModel(p.Weapon.GetWeaponType()))
	default:
		return model.PartProperties{}, nil
	}
}

// EngineClassToModel конвертирует proto-класс двигателя в доменный.
func EngineClassToModel(c inventoryv1.EngineClass) model.EngineClass {
	switch c {
	case inventoryv1.EngineClass_ENGINE_CLASS_A:
		return model.EngineClassA
	case inventoryv1.EngineClass_ENGINE_CLASS_B:
		return model.EngineClassB
	case inventoryv1.EngineClass_ENGINE_CLASS_C:
		return model.EngineClassC
	default:
		return ""
	}
}

// ShieldTypeToModel конвертирует proto-тип щита в доменный.
func ShieldTypeToModel(t inventoryv1.ShieldType) model.ShieldType {
	switch t {
	case inventoryv1.ShieldType_SHIELD_TYPE_ENERGY:
		return model.ShieldTypeEnergy
	case inventoryv1.ShieldType_SHIELD_TYPE_PLASMA:
		return model.ShieldTypePlasma
	default:
		return ""
	}
}

// WeaponTypeToModel конвертирует proto-тип оружия в доменный.
func WeaponTypeToModel(t inventoryv1.WeaponType) model.WeaponType {
	switch t {
	case inventoryv1.WeaponType_WEAPON_TYPE_LASER:
		return model.WeaponTypeLaser
	case inventoryv1.WeaponType_WEAPON_TYPE_MISSILE:
		return model.WeaponTypeMissile
	default:
		return ""
	}
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
			Properties:    PartPropertiesToDTO(p.Properties()),
		},
		CreatedAt: timestamppb.New(p.CreatedAt()),
	}
}

// PartPropertiesToDTO конвертирует типоспецифичные свойства детали в proto-представление.
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
