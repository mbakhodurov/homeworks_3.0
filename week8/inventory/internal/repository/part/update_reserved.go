package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week8/inventory/internal/model"
)

// UpdateReservedBatch обновляет reserved для нескольких деталей одним запросом.
//
// unnest разворачивает два массива (UUID и reserved) в виртуальную таблицу batch(uuid, reserved) —
// пары элементов на одной позиции становятся одной строкой. UPDATE ... FROM batch джойнит
// её с parts по uuid и обновляет reserved — один round-trip к БД вместо N отдельных UPDATE.
func (r *repository) UpdateReservedBatch(ctx context.Context, parts []model.Part) error {
	const query = `
		UPDATE parts AS p
		SET reserved = batch.reserved
		FROM unnest($1::uuid[], $2::int[]) AS batch(uuid, reserved)
		WHERE p.uuid = batch.uuid
	`

	uuids := make([]string, len(parts))
	reserved := make([]int, len(parts))

	for i, p := range parts {
		uuids[i] = p.UUID().String()
		reserved[i] = p.Reserved()
	}

	if _, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, uuids, reserved); err != nil {
		return fmt.Errorf("обновить резерв: %w", err)
	}

	return nil
}
