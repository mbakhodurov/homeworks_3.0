package part

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/repository/record"
)

// List возвращает детали, отфильтрованные по filter.UUIDs (приоритет) или по filter.PartType.
func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	var (
		rows pgx.Rows
		err  error
	)

	switch {
	case len(filter.UUIDs) > 0:
		query := `SELECT ` + partColumns + ` FROM parts WHERE uuid = ANY($1) ORDER BY name`
		rows, err = r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, filter.UUIDs)
	case filter.PartType != "" && filter.PartType != model.PartTypeUnspecified:
		query := `SELECT ` + partColumns + ` FROM parts WHERE part_type = $1 ORDER BY name`
		rows, err = r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, string(filter.PartType))
	default:
		query := `SELECT ` + partColumns + ` FROM parts ORDER BY name`
		rows, err = r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	parts := make([]model.Part, 0, len(recs))
	for _, rec := range recs {
		p, convErr := converter.PartRecordToModel(rec)
		if convErr != nil {
			return nil, fmt.Errorf("конвертировать деталь: %w", convErr)
		}
		parts = append(parts, p)
	}

	return parts, nil
}
