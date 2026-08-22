package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/iam/internal/api/converter"
	userv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/user/v1"
)

// GetUser возвращает пользователя по UUID.
func (a *api) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := a.userService.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: converter.UserToDTO(user),
	}, nil
}
