package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week8/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/service/input"
)

// UserService определяет контракт бизнес-логики управления пользователями.
type UserService interface {
	Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID string) (model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}
