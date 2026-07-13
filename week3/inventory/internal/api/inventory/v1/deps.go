package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/model"
)

// PartService определяет контракт для бизнес-логики работы с деталями.
type PartService interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, uuids []uuid.UUID, partType model.PartType) ([]model.Part, error)
	Create(ctx context.Context, part model.Part) (uuid.UUID, error)
	Delete(ctx context.Context, partUUID uuid.UUID) error
	GetAll(ctx context.Context) ([]model.Part, error)
}
