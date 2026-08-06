package part

import (
	"context"
	"fmt"

	errs "github.com/mbakhodurov/homeworks2/week6/inventory/internal/errors"
)

// ReleaseParts освобождает ранее зарезервированные детали по списку UUID в транзакции:
// читает детали → освобождает каждую через Release() → сохраняет батчем.
func (s *service) ReleaseParts(ctx context.Context, uuids []string) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := s.partRepo.ListForUpdate(ctx, uuids)
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
