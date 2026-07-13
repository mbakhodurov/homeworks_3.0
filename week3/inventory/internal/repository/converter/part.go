package converter

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/record"
)

// PartToRecord конвертирует доменную модель в persistence-shape.
func PartToRecord(p model.Part) record.Part {
	return record.Part{
		UUID:          p.UUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      string(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}

// PartToModel конвертирует запись хранилища в доменную модель.
func PartToModel(r record.Part) model.Part {
	return model.Part{
		UUID:          uuid.MustParse(r.UUID),
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		PartType:      model.PartType(r.PartType),
		StockQuantity: r.StockQuantity,
		CreatedAt:     r.CreatedAt,
	}
}
