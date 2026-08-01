package assembly_consumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/mbakhodurov/homeworks2/week5/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week5/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week5/platform/pkg/kafka"
)

// handler обрабатывает событие ShipAssembled из Kafka:
// при успехе списывает детали и переводит заказ в статус ASSEMBLED.
type handler struct {
	orderRepo       OrderRepository
	inventoryClient InventoryClient
	txManager       TxManager
}

func newHandler(orderRepo OrderRepository, inventoryClient InventoryClient, txManager TxManager) *handler {
	return &handler{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		txManager:       txManager,
	}
}

// Handle обрабатывает одно Kafka-сообщение.
// При ошибке десериализации логирует и возвращает nil (оффсет коммитится — повтор бессмысленен).
// При прочих ошибках возвращает error (оффсет НЕ коммитится, сообщение будет повторено).
func (h *handler) Handle(ctx context.Context, msg kafka.Message) error {
	event, err := decodeShipAssembled(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось десериализовать ShipAssembled — пропускаем", "error", err)
		return nil
	}

	slog.InfoContext(ctx, "обрабатываем ShipAssembled",
		"order_uuid", event.OrderUUID,
		"build_time_sec", event.BuildTimeSec,
	)

	return h.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := h.orderRepo.Get(txCtx, event.OrderUUID)
		if err != nil {
			return fmt.Errorf("получить заказ: %w", err)
		}

		if order.Status == model.OrderStatusAssembled {
			slog.InfoContext(ctx, "заказ уже собран, пропускаем (идемпотентность)", "order_uuid", event.OrderUUID)
			return nil
		}

		if order.Status != model.OrderStatusPaid {
			slog.WarnContext(ctx, "неожиданный статус заказа при сборке",
				"order_uuid", event.OrderUUID,
				"status", order.Status,
			)
			return errs.ErrOrderNotFound
		}

		partUUIDs := make([]uuid.UUID, 0, len(order.Items))
		for _, item := range order.Items {
			partUUIDs = append(partUUIDs, item.PartUUID)
		}

		if err := h.inventoryClient.CommitParts(txCtx, partUUIDs); err != nil {
			return fmt.Errorf("списать детали: %w", err)
		}

		order.Status = model.OrderStatusAssembled
		order.UpdatedAt = lo.ToPtr(time.Now())

		if err := h.orderRepo.Update(txCtx, order); err != nil {
			return fmt.Errorf("обновить заказ: %w", err)
		}

		slog.InfoContext(ctx, "заказ переведён в статус ASSEMBLED", "order_uuid", event.OrderUUID)
		return nil
	})
}
