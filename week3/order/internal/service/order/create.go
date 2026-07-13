package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week3/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
)

// Create создаёт новый заказ атомарно: INSERT в orders + INSERT в order_items в одной транзакции.
// Внешние gRPC-вызовы (InventoryClient) выполняются ДО транзакции.
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
	order.Status = model.OrderStatusPendingPayment
	order.CreatedAt = time.Now()

	for i := range items {
		items[i].UUID = uuid.New()
		items[i].OrderUUID = order.UUID
	}
	order.Items = items

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		if err := s.orderRepo.Create(ctx, order); err != nil {
			return fmt.Errorf("создать заказ: %w", err)
		}
		if err := s.orderItemRepo.Create(ctx, order.Items); err != nil {
			return fmt.Errorf("создать позиции заказа: %w", err)
		}
		return nil
	})
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
