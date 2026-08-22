package part

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week8/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/inventory/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week8/inventory/internal/repository/record"
)

// ListForUpdate возвращает детали с блокировкой строк (SELECT FOR UPDATE) до конца транзакции.
// ORDER BY uuid предотвращает дедлоки при параллельных транзакциях над пересекающимися наборами деталей.
func (r *repository) ListForUpdate(ctx context.Context, uuids []string) ([]model.Part, error) {
	query := `SELECT ` + partColumns + ` FROM parts WHERE uuid = ANY($1) ORDER BY uuid FOR UPDATE`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, uuids)
	if err != nil {
		return nil, fmt.Errorf("получить детали с блокировкой: %w", err)
	}

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить детали с блокировкой: %w", err)
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
