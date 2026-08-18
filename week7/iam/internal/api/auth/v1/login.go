package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week7/iam/internal/service/input"
	authv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/auth/v1"
)

// Login выполняет вход в систему и создаёт сессию.
func (a *api) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	sessionUUID, err := a.authService.Login(ctx, input.LoginInput{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		SessionUuid: sessionUUID.String(),
	}, nil
}
