package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
)

// List возвращает детали, отфильтрованные по uuids (приоритет) или по типу.
func (s *service) List(ctx context.Context, uuids []uuid.UUID, partType model.PartType) ([]model.Part, error) {
	parts, err := s.partRepo.List(ctx, uuids, partType)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return parts, nil
}
