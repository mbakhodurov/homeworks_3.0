package part

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/record"
)

// GetAll возвращает все детали без фильтрации, отсортированные по имени.
func (r *repository) GetAll(ctx context.Context) ([]model.Part, error) {
	const query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts ORDER BY name`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("получить все детали: %w", err)
	}

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить все детали: %w", err)
	}

	parts := make([]model.Part, 0, len(recs))
	for _, rec := range recs {
		parts = append(parts, converter.PartToModel(rec))
	}

	return parts, nil
}
