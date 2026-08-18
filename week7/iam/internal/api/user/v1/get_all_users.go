package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week7/iam/internal/api/converter"
	commonv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/common/v1"
	userv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/user/v1"
)

func (a *api) GetAllUsers(ctx context.Context, _ *userv1.GetAllUsersRequest) (*userv1.GetAllUsersResponse, error) {
	users, err := a.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*commonv1.User, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, converter.UserToDTO(u))
	}

	return &userv1.GetAllUsersResponse{Users: dtos}, nil
}
