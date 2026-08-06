package iam

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/model"
)

// UserRepository определяет контракт для работы с хранилищем пользователей.
type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByUUID(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}

// SessionRepository определяет контракт для работы с хранилищем сессий.
type SessionRepository interface {
	Create(ctx context.Context, session model.Session, ttl time.Duration) error
	Get(ctx context.Context, sessionUUID string) (model.Session, error)
	Delete(ctx context.Context, sessionUUID string) error
}
