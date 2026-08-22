package orderitem

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/mbakhodurov/homeworks2/week8/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/order/internal/repository/converter"
)

// Create сохраняет позиции заказа: INSERT в order_items.
// Вызывается внутри транзакции сервисного слоя.
func (r *repository) Create(ctx context.Context, items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	itemRecs := converter.OrderItemsToRecords(items)

	q := squirrel.Insert("order_items").
		Columns("uuid", "order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range itemRecs {
		q = q.Values(item.UUID, item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("собрать запрос позиций заказа: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("создать позиции заказа: %w", err)
	}

	return nil
}
