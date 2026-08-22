package converter

import (
	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week8/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/repository/record"
)

// UserToRecord конвертирует доменную модель пользователя в persistence-shape.
func UserToRecord(u model.User) record.User {
	return record.User{
		UUID:         u.UUID.String(),
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// UserFromRecord восстанавливает доменную модель пользователя из persistence-shape.
func UserFromRecord(r record.User) model.User {
	return model.User{
		UUID:         uuid.MustParse(r.UUID),
		Login:        r.Login,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
