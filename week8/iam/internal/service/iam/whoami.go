package iam

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week8/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/model"
)

// Whoami проверяет сессию и возвращает её саму вместе с владельцем.
func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	if sessionUUID == "" {
		return model.Session{}, model.User{}, errs.ErrEmptySessionID
	}

	session, err := s.sessionRepo.Get(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, errs.ErrSessionNotFound) {
			return model.Session{}, model.User{}, errs.ErrSessionNotFound
		}

		return model.Session{}, model.User{}, fmt.Errorf("проверить сессию: %w", err)
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return model.Session{}, model.User{}, fmt.Errorf("получить пользователя сессии: %w", err)
	}

	return session, user, nil
}
