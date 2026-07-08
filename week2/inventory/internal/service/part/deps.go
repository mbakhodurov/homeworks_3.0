package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
)

// PartRepository определяет контракт для работы с хранилищем деталей.
type PartRepository interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, uuids []uuid.UUID, partType model.PartType) ([]model.Part, error)
	Create(ctx context.Context, part model.Part) error
	Delete(ctx context.Context, partUUID uuid.UUID) error
	GetAll(ctx context.Context) ([]model.Part, error)
}
