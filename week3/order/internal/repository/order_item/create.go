package orderitem

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
)

func (r *repository) Create(ctx context.Context, items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	query := squirrel.Insert("order_items").
		Columns("uuid", "order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range items {
		query = query.Values(item.UUID, item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("insert order items: %w", err)
	}

	return nil
}
