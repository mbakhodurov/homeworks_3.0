package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/mbakhodurov/homeworks2/week8/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/service/input"
)

// minPasswordLength — минимальная длина пароля при регистрации.
const minPasswordLength = 8

// Register регистрирует нового пользователя и возвращает его UUID.
func (s *service) Register(ctx context.Context, in input.RegisterInput) (uuid.UUID, error) {
	if in.Login == "" {
		return uuid.Nil, errs.ErrInvalidLogin
	}

	if len(in.Password) < minPasswordLength {
		return uuid.Nil, errs.ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), s.bcryptCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("захешировать пароль: %w", err)
	}

	user := model.User{
		UUID:         uuid.New(),
		Login:        in.Login,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return uuid.Nil, fmt.Errorf("зарегистрировать пользователя: %w", err)
	}

	return user.UUID, nil
}
