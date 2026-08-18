package iam

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/mbakhodurov/homeworks2/week7/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/service/input"
)

// Login проверяет логин/пароль, создаёт сессию и возвращает её UUID.
//
// И «логин не найден», и «неверный пароль» возвращают одну и ту же ошибку
// ErrInvalidCredentials — так атакующий не может понять по ответу API,
// существует ли логин в системе (частичная защита от username enumeration).
func (s *service) Login(ctx context.Context, in input.LoginInput) (uuid.UUID, error) {
	if in.Login == "" || in.Password == "" {
		return uuid.Nil, errs.ErrEmptyCredential
	}

	user, err := s.userRepo.GetByLogin(ctx, in.Login)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return uuid.Nil, errs.ErrInvalidCredentials
		}

		return uuid.Nil, fmt.Errorf("войти: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return uuid.Nil, errs.ErrInvalidCredentials
	}

	now := time.Now()
	session := model.Session{
		UUID:      uuid.New(),
		UserUUID:  user.UUID,
		Login:     user.Login,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}

	if err = s.sessionRepo.Create(ctx, session, s.sessionTTL); err != nil {
		return uuid.Nil, fmt.Errorf("создать сессию: %w", err)
	}

	return session.UUID, nil
}
