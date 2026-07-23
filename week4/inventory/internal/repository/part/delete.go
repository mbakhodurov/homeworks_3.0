package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
)

// Delete удаляет деталь по UUID.
func (r *repository) Delete(ctx context.Context, partUUID uuid.UUID) error {
	const query = `DELETE FROM parts WHERE uuid = $1`

	result, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, partUUID)
	if err != nil {
		return fmt.Errorf("удалить деталь: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrPartNotFound
	}

	return nil
}
