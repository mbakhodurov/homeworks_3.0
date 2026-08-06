package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week6/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week6/iam/internal/model"
)

// GetUser возвращает пользователя по UUID.
func (s *service) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	id, err := uuid.Parse(userUUID)
	if err != nil {
		return model.User{}, errs.ErrInvalidUUID
	}

	user, err := s.userRepo.GetByUUID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return model.User{}, errs.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("получить пользователя: %w", err)
	}

	return user, nil
}
