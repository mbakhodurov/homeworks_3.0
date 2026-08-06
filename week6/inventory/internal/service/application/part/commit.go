package part

import (
	"context"
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week6/inventory/internal/errors"
)

// CommitParts списывает зарезервированные детали со склада после сборки корабля.
func (s *service) CommitParts(ctx context.Context, uuids []string) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := s.partRepo.ListForUpdate(ctx, uuids)
		if err != nil {
			return fmt.Errorf("получить детали для списания: %w", err)
		}

		if len(parts) != len(uuids) {
			return errs.ErrPartNotFound
		}

		for _, p := range parts {
			if p.StockQuantity() <= 0 || p.Reserved() <= 0 {
				return errs.ErrNothingToCommit
			}
		}

		if err = s.partRepo.CommitParts(ctx, uuids); err != nil {
			return fmt.Errorf("списать детали: %w", err)
		}

		return nil
	})
}
