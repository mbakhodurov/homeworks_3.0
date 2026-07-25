package model

import (
	"time"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
)

// Part — доменная сущность детали космического корабля.
// Инварианты: reserved не может превысить stockQuantity и не может стать отрицательным.
type Part struct {
	uuid          uuid.UUID
	name          string
	description   string
	partType      PartType
	price         int64
	stockQuantity int
	reserved      int
	properties    PartProperties
	createdAt     time.Time
}

// RestorePart восстанавливает сущность из БД (без валидации — данные уже проверены).
func RestorePart(
	partUUID uuid.UUID,
	name, description string,
	partType PartType,
	price int64,
	stockQuantity, reserved int,
	properties PartProperties,
	createdAt time.Time,
) Part {
	return Part{
		uuid:          partUUID,
		name:          name,
		description:   description,
		partType:      partType,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		properties:    properties,
		createdAt:     createdAt,
	}
}

// Reserve резервирует одну единицу детали.
// Возвращает ErrOutOfStock, если свободных единиц нет.
func (p *Part) Reserve() error {
	if p.stockQuantity-p.reserved <= 0 {
		return errs.ErrOutOfStock
	}

	p.reserved++

	return nil
}

// Release освобождает одну зарезервированную единицу детали.
// Возвращает ErrNothingToRelease, если резерв равен нулю.
func (p *Part) Release() error {
	if p.reserved <= 0 {
		return errs.ErrNothingToRelease
	}

	p.reserved--

	return nil
}

// UUID возвращает уникальный идентификатор детали.
func (p *Part) UUID() uuid.UUID { return p.uuid }

// Name возвращает название детали.
func (p *Part) Name() string { return p.name }

// Description возвращает описание детали.
func (p *Part) Description() string { return p.description }

// PartType возвращает тип детали.
func (p *Part) PartType() PartType { return p.partType }

// Price возвращает цену детали в копейках.
func (p *Part) Price() int64 { return p.price }

// StockQuantity возвращает общее количество на складе.
func (p *Part) StockQuantity() int { return p.stockQuantity }

// Reserved возвращает количество зарезервированных единиц.
func (p *Part) Reserved() int { return p.reserved }

// Properties возвращает типоспецифичные свойства детали.
func (p *Part) Properties() *PartProperties { return &p.properties }

// CreatedAt возвращает время создания записи.
func (p *Part) CreatedAt() time.Time { return p.createdAt }

// PartFilter задаёт параметры фильтрации деталей.
// Если UUIDs непустой — фильтрация по UUID. Иначе — по PartType
// (или все детали, если PartType == "" / PartTypeUnspecified).
type PartFilter struct {
	UUIDs    []string
	PartType PartType
}

// CreatePartInput задаёт параметры создания новой детали.
type CreatePartInput struct {
	Name          string
	Description   string
	PartType      PartType
	Price         int64
	StockQuantity int
	Properties    PartProperties
}
