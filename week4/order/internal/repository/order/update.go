package order

import (
	"context"
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week4/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/order/internal/repository/converter"
)

// Update сохраняет изменения существующего заказа (статус, транзакция, способ оплаты, updated_at).
func (r *repository) Update(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	const query = `
		UPDATE orders
		SET status = $1, transaction_uuid = $2, payment_method = $3, updated_at = $4
		WHERE uuid = $5`

	tag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, rec.Status, rec.TransactionUUID, rec.PaymentMethod, rec.UpdatedAt, rec.UUID)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}
