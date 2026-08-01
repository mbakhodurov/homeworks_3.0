package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week5/inventory/internal/model"
)

// TxManager определяет контракт для управления транзакциями.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// CompatibilityChecker определяет контракт для доменного сервиса проверки совместимости.
type CompatibilityChecker interface {
	Check(parts []model.Part) error
}

// PartRepository определяет контракт для работы с хранилищем деталей.
type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
	ListForUpdate(ctx context.Context, uuids []string) ([]model.Part, error)
	Create(ctx context.Context, p model.Part) error
	Delete(ctx context.Context, partUUID uuid.UUID) error
	UpdateReservedBatch(ctx context.Context, parts []model.Part) error
	CommitParts(ctx context.Context, uuids []string) error
}
