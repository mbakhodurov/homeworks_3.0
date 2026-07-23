package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week4/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/input"
)

// ReleaseParts освобождает ранее зарезервированные детали по списку UUID в транзакции:
// читает детали → освобождает каждую через Release() → сохраняет батчем.
func (s *service) ReleaseParts(ctx context.Context, uuids []uuid.UUID) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := s.partRepo.List(ctx, input.PartFilter{UUIDs: uuids})
		if err != nil {
			return fmt.Errorf("получить детали: %w", err)
		}

		if len(parts) != len(uuids) {
			return errs.ErrPartNotFound
		}

		for i := range parts {
			if err = parts[i].Release(); err != nil {
				return fmt.Errorf("освободить деталь %s: %w", parts[i].Name(), err)
			}
		}

		if err = s.partRepo.UpdateReservedBatch(ctx, parts); err != nil {
			return fmt.Errorf("сохранить резерв: %w", err)
		}

		return nil
	})
}
