package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
)

// PartService определяет контракт для бизнес-логики работы с деталями.
type PartService interface {
	Get(ctx context.Context, partUUID string) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
	Create(ctx context.Context, in model.CreatePartInput) (uuid.UUID, error)
	Delete(ctx context.Context, partUUID uuid.UUID) error
	ValidateCompatibility(ctx context.Context, uuids []string) error
	ReserveParts(ctx context.Context, uuids []string) error
	ReleaseParts(ctx context.Context, uuids []string) error
}
