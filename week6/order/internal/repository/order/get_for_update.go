package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/mbakhodurov/homeworks2/week6/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week6/order/internal/repository/record"
)

// GetForUpdate возвращает заказ по UUID с блокировкой строки (SELECT FOR UPDATE).
// Используется в сценариях, где параллельные запросы должны ждать завершения транзакции.
func (r *repository) GetForUpdate(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	const orderQuery = `
        SELECT uuid, user_uuid, total_price, status, transaction_uuid, payment_method, created_at, updated_at
        FROM orders
        WHERE uuid = $1
        FOR UPDATE`

	orderRows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, orderQuery, orderUUID.String())
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ с блокировкой: %w", err)
	}

	orderRec, err := pgx.CollectOneRow(orderRows, pgx.RowToStructByName[record.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, errs.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("получить заказ с блокировкой: %w", err)
	}

	const itemsQuery = `
        SELECT uuid, order_uuid, part_uuid, part_type, price, created_at
        FROM order_items
        WHERE order_uuid = $1`

	itemRows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, itemsQuery, orderUUID.String())
	if err != nil {
		return model.Order{}, fmt.Errorf("получить позиции заказа: %w", err)
	}

	itemRecs, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return model.Order{}, fmt.Errorf("получить позиции заказа: %w", err)
	}

	return converter.OrderToModel(orderRec, itemRecs), nil
}
