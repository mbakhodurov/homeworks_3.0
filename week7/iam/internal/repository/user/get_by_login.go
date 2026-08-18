package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errs "github.com/mbakhodurov/homeworks2/week7/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/repository/record"
)

// GetByLogin возвращает пользователя по логину.
func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	const q = `
		SELECT uuid, login, password_hash, created_at, updated_at
		FROM users
		WHERE login = $1`

	rows, err := r.pool.Query(ctx, q, login)
	if err != nil {
		return model.User{}, fmt.Errorf("получить пользователя по логину: %w", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errs.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("получить пользователя по логину: %w", err)
	}

	return converter.UserFromRecord(rec), nil
}
