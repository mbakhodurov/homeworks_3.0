package v1

import (
	userv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/user/v1"
)

// api реализует gRPC-обработчики UserService.
type api struct {
	userv1.UnimplementedUserServiceServer

	userService UserService
}

// New создаёт API-обработчик UserService.
func New(userService UserService) *api {
	return &api{
		userService: userService,
	}
}
