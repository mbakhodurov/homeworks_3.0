package part

import (
	"context"
	"fmt"
)

// CommitParts списывает stock_quantity и reserved на 1 для каждой детали из uuids.
// Условие stock_quantity > 0 AND reserved > 0 исключает уже списанные детали.
func (r *repository) CommitParts(ctx context.Context, uuids []string) error {
	const query = `
		UPDATE parts
		SET stock_quantity = stock_quantity - 1,
		    reserved       = reserved - 1
		WHERE uuid = ANY($1)
		  AND stock_quantity > 0
		  AND reserved > 0
	`

	tag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, uuids)
	if err != nil {
		return fmt.Errorf("списать детали: %w", err)
	}

	if tag.RowsAffected() != int64(len(uuids)) {
		return fmt.Errorf("списать детали: обновлено %d из %d", tag.RowsAffected(), len(uuids))
	}

	return nil
}
