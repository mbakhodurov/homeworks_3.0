package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Delete удаляет деталь по UUID.
func (s *service) Delete(ctx context.Context, partUUID uuid.UUID) error {
	if err := s.partRepo.Delete(ctx, partUUID); err != nil {
		return fmt.Errorf("удалить деталь: %w", err)
	}

	return nil
}
