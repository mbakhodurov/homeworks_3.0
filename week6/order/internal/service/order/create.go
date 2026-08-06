package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week6/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/auth"
)

// Create создаёт новый заказ.
//
// Алгоритм:
//  1. Получает детали по UUID слотов (ListParts) — до транзакции.
//  2. Проверяет совместимость деталей (ValidateCompatibility) — до транзакции.
//  3. Резервирует детали (ReserveParts) — до транзакции.
//  4. Атомарно сохраняет заказ и его позиции: INSERT в orders + INSERT в order_items.
//
// Внешние gRPC-вызовы к InventoryService выполняются ДО транзакции — так дольше,
// чем длится сама транзакция, не держится открытое соединение к БД.
func (s *service) Create(ctx context.Context, in model.CreateOrderInput) (model.Order, error) {
	userUUID, ok := auth.UserUUIDFromContext(ctx)
	if !ok {
		return model.Order{}, errs.ErrUnauthorized
	}

	uuids := in.PartUUIDs()

	parts, err := s.inventoryClient.ListParts(ctx, uuids)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали: %w", err)
	}

	slots := model.ShipSlots(in)

	if err = s.inventoryClient.ValidateCompatibility(ctx, slots); err != nil {
		return model.Order{}, fmt.Errorf("проверить совместимость деталей: %w", err)
	}

	if err = s.inventoryClient.ReserveParts(ctx, uuids); err != nil {
		return model.Order{}, fmt.Errorf("зарезервировать детали: %w", err)
	}

	order := model.Order{
		UUID:      uuid.New(),
		UserUUID:  userUUID,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	items := make([]model.OrderItem, 0, len(parts))
	for _, part := range parts {
		items = append(items, model.OrderItem{
			UUID:      uuid.New(),
			OrderUUID: order.UUID,
			PartUUID:  part.UUID,
			PartType:  part.PartType,
			Price:     part.Price,
		})
	}
	order.Items = items

	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := s.orderRepo.Create(txCtx, order); err != nil {
			return fmt.Errorf("создать заказ: %w", err)
		}
		if err := s.orderItemRepo.Create(txCtx, items); err != nil {
			return fmt.Errorf("создать позиции заказа: %w", err)
		}
		return nil
	})
	if err != nil {
		// Резерв в InventoryService уже создан, а заказ в своей БД — нет.
		// Это две независимые БД без общей транзакции, поэтому откатываем
		// резерв компенсирующим вызовом (best-effort: полноценная сага —
		// тема будущих недель). Если откат тоже не удался — деталь
		// останется зависшей в резерве, логируем для ручного разбора.
		if releaseErr := s.inventoryClient.ReleaseParts(ctx, uuids); releaseErr != nil {
			slog.ErrorContext(ctx, "не удалось откатить резерв деталей после ошибки создания заказа",
				"part_uuids", uuids, "create_error", err, "release_error", releaseErr)
		}

		return model.Order{}, err
	}

	return order, nil
}
