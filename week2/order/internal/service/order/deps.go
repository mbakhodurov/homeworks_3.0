package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
)

// OrderRepository определяет контракт для работы с хранилищем заказов.
type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
	GetAll(ctx context.Context) ([]model.Order, error)
}

// InventoryClient определяет контракт для работы с InventoryService.
type InventoryClient interface {
	ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error)
}

// PaymentClient определяет контракт для работы с PaymentService.
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
}
