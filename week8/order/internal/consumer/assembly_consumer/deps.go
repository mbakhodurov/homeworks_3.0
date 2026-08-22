package assembly_consumer

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week8/order/internal/model"
)

// OrderRepository определяет контракт для получения и обновления заказов.
type OrderRepository interface {
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

// InventoryClient определяет контракт для списания деталей.
type InventoryClient interface {
	CommitParts(ctx context.Context, uuids []uuid.UUID) error
}

// TxManager определяет контракт для управления транзакциями.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
