package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/record"
)

// List возвращает детали, отфильтрованные по uuids (приоритет) или по типу.
func (r *repository) List(ctx context.Context, uuids []uuid.UUID, partType model.PartType) ([]model.Part, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if len(uuids) > 0 {
		uuidStrings := make([]string, len(uuids))
		for i, u := range uuids {
			uuidStrings[i] = u.String()
		}

		const query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts WHERE uuid = ANY($1) ORDER BY name`
		rows, err = r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, uuidStrings)
	} else if partType != model.PartTypeUnspecified {
		const query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts WHERE part_type = $1 ORDER BY name`
		rows, err = r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, string(partType))
	} else {
		const query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts ORDER BY name`
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
		parts = append(parts, converter.PartToModel(rec))
	}

	return parts, nil
}
