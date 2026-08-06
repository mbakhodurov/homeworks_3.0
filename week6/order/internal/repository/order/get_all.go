package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/repository/record"
)

// GetAll возвращает все заказы вместе с позициями.
func (r *repository) GetAll(ctx context.Context) ([]model.Order, error) {
	const ordersQuery = `
		SELECT uuid, user_uuid, total_price, status, transaction_uuid, payment_method, created_at, updated_at
		FROM orders`

	orderRows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, ordersQuery)
	if err != nil {
		return nil, fmt.Errorf("получить заказы: %w", err)
	}

	orderRecs, err := pgx.CollectRows(orderRows, pgx.RowToStructByName[record.Order])
	if err != nil {
		return nil, fmt.Errorf("получить заказы: %w", err)
	}

	if len(orderRecs) == 0 {
		return nil, nil
	}

	uuids := make([]string, 0, len(orderRecs))
	for _, o := range orderRecs {
		uuids = append(uuids, o.UUID)
	}

	const itemsQuery = `
		SELECT uuid, order_uuid, part_uuid, part_type, price, created_at
		FROM order_items
		WHERE order_uuid = ANY($1)`

	itemRows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, itemsQuery, uuids)
	if err != nil {
		return nil, fmt.Errorf("получить позиции заказов: %w", err)
	}

	itemRecs, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("получить позиции заказов: %w", err)
	}

	itemsByOrder := make(map[string][]record.OrderItem, len(orderRecs))
	for _, item := range itemRecs {
		itemsByOrder[item.OrderUUID] = append(itemsByOrder[item.OrderUUID], item)
	}

	orders := make([]model.Order, 0, len(orderRecs))
	for _, o := range orderRecs {
		orders = append(orders, converter.OrderToModel(o, itemsByOrder[o.UUID]))
	}

	return orders, nil
}
