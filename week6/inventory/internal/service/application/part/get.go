package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/model"
)

// Get возвращает деталь по UUID.
func (s *service) Get(ctx context.Context, partUUID string) (model.Part, error) {
	p, err := s.partRepo.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return p, nil
}
