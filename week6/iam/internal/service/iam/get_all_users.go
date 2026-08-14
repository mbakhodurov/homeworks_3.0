package iam

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/model"
)

func (s *service) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("получить всех пользователей: %w", err)
	}

	return users, nil
}
