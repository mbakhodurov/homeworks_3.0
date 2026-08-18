package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mbakhodurov/homeworks2/week7/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/repository/record"
)

func (r *repository) GetAll(ctx context.Context) ([]model.User, error) {
	const q = `
		SELECT uuid, login, password_hash, created_at, updated_at
		FROM users`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("получить всех пользователей: %w", err)
	}

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.User])
	if err != nil {
		return nil, fmt.Errorf("получить всех пользователей: %w", err)
	}

	users := make([]model.User, len(recs))
	for i, rec := range recs {
		users[i] = converter.UserFromRecord(rec)
	}

	return users, nil
}
