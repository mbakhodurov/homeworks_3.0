package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
)

// GetAll возвращает все детали без фильтрации.
func (s *service) GetAll(ctx context.Context) ([]model.Part, error) {
	parts, err := s.partRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("получить все детали: %w", err)
	}

	return parts, nil
}
