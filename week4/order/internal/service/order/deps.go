package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/order/internal/model"
)

// OrderRepository определяет контракт для работы с хранилищем заказов.
type OrderRepository interface {
	Create(ctx context.Context, order model.Order, items []model.OrderItem) error
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

// InventoryClient определяет контракт для работы с InventoryService.
type InventoryClient interface {
	ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error)
	ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error
	ReserveParts(ctx context.Context, uuids []uuid.UUID) error
	ReleaseParts(ctx context.Context, uuids []uuid.UUID) error
}

// PaymentClient определяет контракт для работы с PaymentService.
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
}

// TxManager определяет контракт для управления транзакциями.
// Реализация — *manager.Manager из библиотеки go-transaction-manager.
//
// Метод Do оборачивает callback в SQL-транзакцию (BEGIN/COMMIT/ROLLBACK):
//   - Перед вызовом fn: выполняет BEGIN и кладёт pgx.Tx внутрь нового ctx
//   - fn возвращает nil → выполняет COMMIT
//   - fn возвращает error → выполняет ROLLBACK
//
// Ключевой момент: ctx, который получает fn, содержит внутри себя активную
// транзакцию. Именно этот ctx нужно передавать дальше в репозиторий — тогда
// все запросы пойдут через одно соединение в рамках одной транзакции.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
