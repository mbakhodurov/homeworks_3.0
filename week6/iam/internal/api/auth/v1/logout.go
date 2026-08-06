package v1

import (
	"context"

	authv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/auth/v1"
)

// Logout завершает сессию пользователя.
func (a *api) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := a.authService.Logout(ctx, req.GetSessionUuid()); err != nil {
		return nil, err
	}

	return &authv1.LogoutResponse{}, nil
}
