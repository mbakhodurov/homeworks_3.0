package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/api/converter"
	authv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/auth/v1"
)

// Whoami проверяет текущую сессию и возвращает информацию о пользователе.
func (a *api) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	session, user, err := a.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &authv1.WhoamiResponse{
		Session: converter.SessionToDTO(session),
		User:    converter.UserToDTO(user),
	}, nil
}
