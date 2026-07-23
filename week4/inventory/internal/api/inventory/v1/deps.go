package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/input"
)

// PartService определяет контракт для бизнес-логики работы с деталями.
type PartService interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
	Create(ctx context.Context, in input.CreatePartInput) (uuid.UUID, error)
	Delete(ctx context.Context, partUUID uuid.UUID) error
	ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error
	ReserveParts(ctx context.Context, uuids []uuid.UUID) error
	ReleaseParts(ctx context.Context, uuids []uuid.UUID) error
}
