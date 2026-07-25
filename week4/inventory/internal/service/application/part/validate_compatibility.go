package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
)

// ValidateCompatibility проверяет совместимость деталей по переданным UUID.
func (s *service) ValidateCompatibility(ctx context.Context, uuids []string) error {
	parts, err := s.partRepo.List(ctx, model.PartFilter{UUIDs: uuids})
	if err != nil {
		return fmt.Errorf("получить детали: %w", err)
	}

	return s.compatibilityChecker.Check(parts)
}
