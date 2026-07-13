package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week3/order/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week3/order/internal/repository/record"
)

// GetAll возвращает все заказы со всеми позициями.
func (r *repository) GetAll(ctx context.Context) ([]model.Order, error) {
	const ordersQuery = `
		SELECT uuid, total_price, status, transaction_uuid, payment_method, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC`

	orderRows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, ordersQuery)
	if err != nil {
		return nil, fmt.Errorf("получить все заказы: %w", err)
	}

	orderRecs, err := pgx.CollectRows(orderRows, pgx.RowToStructByName[record.Order])
	if err != nil {
		return nil, fmt.Errorf("получить все заказы: %w", err)
	}

	if len(orderRecs) == 0 {
		return []model.Order{}, nil
	}

	const itemsQuery = `
		SELECT uuid, order_uuid, part_uuid, part_type, price, created_at
		FROM order_items`

	itemRows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, itemsQuery)
	if err != nil {
		return nil, fmt.Errorf("получить все позиции заказов: %w", err)
	}

	itemRecs, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("получить все позиции заказов: %w", err)
	}

	itemsByOrder := make(map[string][]record.OrderItem, len(orderRecs))
	for _, item := range itemRecs {
		itemsByOrder[item.OrderUUID] = append(itemsByOrder[item.OrderUUID], item)
	}

	orders := make([]model.Order, 0, len(orderRecs))
	for _, orderRec := range orderRecs {
		orders = append(orders, converter.OrderToModel(orderRec, itemsByOrder[orderRec.UUID]))
	}

	return orders, nil
}
