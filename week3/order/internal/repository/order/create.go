package order

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week3/order/internal/repository/converter"
)

// Create сохраняет заказ в таблицу orders (без позиций — они сохраняются отдельно).
func (r *repository) Create(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	const query = `
		INSERT INTO orders (uuid, total_price, status, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, rec.UUID, rec.TotalPrice, rec.Status, rec.CreatedAt)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}
