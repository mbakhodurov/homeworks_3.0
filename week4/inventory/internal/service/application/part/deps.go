package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/input"
)

// TxManager определяет контракт для управления транзакциями.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// CompatibilityChecker определяет контракт для доменного сервиса проверки совместимости.
type CompatibilityChecker interface {
	Check(slots model.ResolvedShipSlots) error
}

// PartRepository определяет контракт для работы с хранилищем деталей.
type PartRepository interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
	Create(ctx context.Context, p model.Part) error
	Delete(ctx context.Context, partUUID uuid.UUID) error
	UpdateReservedBatch(ctx context.Context, parts []model.Part) error
}
