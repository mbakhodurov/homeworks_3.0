package order

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/mbakhodurov/homeworks2/week4/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/order/internal/repository/converter"
)

// Create сохраняет заказ вместе с позициями: INSERT в orders + INSERT в order_items.
// Вызывается внутри транзакции сервисного слоя — оба INSERT атомарны.
func (r *repository) Create(ctx context.Context, order model.Order, items []model.OrderItem) error {
	rec := converter.OrderToRecord(order)

	const orderQuery = `
		INSERT INTO orders (uuid, total_price, status, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, orderQuery, rec.UUID, rec.TotalPrice, rec.Status, rec.CreatedAt)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	if len(items) == 0 {
		return nil
	}

	itemRecs := converter.OrderItemsToRecords(items)

	itemsQuery := squirrel.Insert("order_items").
		Columns("uuid", "order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range itemRecs {
		itemsQuery = itemsQuery.Values(item.UUID, item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := itemsQuery.ToSql()
	if err != nil {
		return fmt.Errorf("собрать запрос позиций заказа: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("создать позиции заказа: %w", err)
	}

	return nil
}
