package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week6/iam/internal/service/input"
)

// AuthService определяет контракт бизнес-логики аутентификации.
type AuthService interface {
	Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID string) error
}
