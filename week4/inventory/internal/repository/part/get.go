package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/repository/record"
)

// Get возвращает деталь по UUID.
func (r *repository) Get(ctx context.Context, partUUID string) (model.Part, error) {
	query := `SELECT ` + partColumns + ` FROM parts WHERE uuid = $1`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, partUUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Part{}, errs.ErrPartNotFound
		}
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	p, err := converter.PartRecordToModel(rec)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать деталь: %w", err)
	}

	return p, nil
}
