package iam

import (
	"context"
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week8/iam/internal/errors"
)

// Logout удаляет сессию. Идемпотентен — повторный Logout для несуществующей
// сессии не считается ошибкой.
func (s *service) Logout(ctx context.Context, sessionUUID string) error {
	if sessionUUID == "" {
		return errs.ErrEmptySessionID
	}

	if err := s.sessionRepo.Delete(ctx, sessionUUID); err != nil {
		return fmt.Errorf("выйти из системы: %w", err)
	}

	return nil
}
