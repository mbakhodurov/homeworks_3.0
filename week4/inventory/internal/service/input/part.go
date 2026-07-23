package input

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
)

// PartFilter задаёт параметры фильтрации деталей.
// Если UUIDs непустой — фильтрация по UUID. Иначе — по PartType
// (или все детали, если PartType == "" / PartTypeUnspecified).
type PartFilter struct {
	UUIDs    []uuid.UUID
	PartType model.PartType
}

// CreatePartInput задаёт параметры создания новой детали.
type CreatePartInput struct {
	Name          string
	Description   string
	PartType      model.PartType
	Price         int64
	StockQuantity int
}
