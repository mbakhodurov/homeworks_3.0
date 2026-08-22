package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/iam/internal/service/input"
	userv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/user/v1"
)

// Register регистрирует нового пользователя.
func (a *api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	userUUID, err := a.userService.Register(ctx, input.RegisterInput{
		Login:    req.GetInfo().GetInfo().GetLogin(),
		Password: req.GetInfo().GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{
		UserUuid: userUUID.String(),
	}, nil
}
