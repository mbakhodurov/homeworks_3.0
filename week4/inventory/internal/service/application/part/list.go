package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
)

// List возвращает детали, отфильтрованные по filter.
func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	parts, err := s.partRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return parts, nil
}
