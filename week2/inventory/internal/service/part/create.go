package part

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
)

// Create создаёт новую деталь и возвращает её UUID.
// UUID и CreatedAt — идентичность записи, их назначает сервис, а не вызывающий код.
func (s *service) Create(ctx context.Context, part model.Part) (uuid.UUID, error) {
	part.UUID = uuid.New()
	part.CreatedAt = time.Now()

	if err := s.partRepo.Create(ctx, part); err != nil {
		return uuid.Nil, fmt.Errorf("создать деталь: %w", err)
	}

	return part.UUID, nil
}
