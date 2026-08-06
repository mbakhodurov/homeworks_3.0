package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/mbakhodurov/homeworks2/week6/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/iam/internal/repository/converter"
)

// uniqueViolationCode — код ошибки PostgreSQL при нарушении unique-constraint.
const uniqueViolationCode = "23505"

// Create сохраняет нового пользователя: INSERT в users.
func (r *repository) Create(ctx context.Context, user model.User) error {
	rec := converter.UserToRecord(user)

	const q = `
		INSERT INTO users (uuid, login, password_hash, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(ctx, q, rec.UUID, rec.Login, rec.PasswordHash, rec.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return errs.ErrUserAlreadyExists
		}

		return fmt.Errorf("создать пользователя: %w", err)
	}

	return nil
}
