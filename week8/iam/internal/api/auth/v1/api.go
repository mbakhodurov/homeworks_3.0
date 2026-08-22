package v1

import (
	authv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/auth/v1"
)

// api реализует gRPC-обработчики AuthService.
type api struct {
	authv1.UnimplementedAuthServiceServer

	authService AuthService
}

// New создаёт API-обработчик AuthService.
func New(authService AuthService) *api {
	return &api{
		authService: authService,
	}
}
