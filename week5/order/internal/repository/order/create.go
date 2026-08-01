package order

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week5/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week5/order/internal/repository/converter"
)

// Create сохраняет заказ: INSERT в orders.
// Вызывается внутри транзакции сервисного слоя.
func (r *repository) Create(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	const q = `
		INSERT INTO orders (uuid, total_price, status, user_uuid, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, q, rec.UUID, rec.TotalPrice, rec.Status, rec.UserUUID, rec.CreatedAt)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}
