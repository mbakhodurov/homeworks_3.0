package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
)

// Get возвращает деталь по UUID.
func (s *service) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	part, err := s.partRepo.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return part, nil
}
