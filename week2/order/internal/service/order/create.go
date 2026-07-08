package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week2/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/order/internal/model"
)

// Create создаёт новый заказ.
// order.Items на входе содержит только запрошенные PartUUID — остальные поля
// (тип, цена) сервис достраивает сам через InventoryClient. UUID, Status и
// CreatedAt — идентичность и жизненный цикл заказа, их назначает сервис.
func (s *service) Create(ctx context.Context, order model.Order) (model.Order, error) {
	uuids := make([]uuid.UUID, 0, len(order.Items))
	for _, item := range order.Items {
		uuids = append(uuids, item.PartUUID)
	}

	parts, err := s.inventoryClient.ListParts(ctx, uuids)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали: %w", err)
	}

	items := make([]model.OrderItem, 0, len(parts))
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		items = append(items, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
	}

	order.UUID = uuid.New()
	order.Items = items
	order.Status = model.OrderStatusPendingPayment
	order.CreatedAt = time.Now()

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return model.Order{}, fmt.Errorf("создать заказ: %w", err)
	}

	return order, nil
}
