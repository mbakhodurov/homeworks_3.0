package part

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
)

// Create создаёт новую деталь и возвращает её UUID.
func (s *service) Create(ctx context.Context, in model.CreatePartInput) (uuid.UUID, error) {
	partUUID := uuid.New()

	newPart := model.RestorePart(
		partUUID,
		in.Name,
		in.Description,
		in.PartType,
		in.Price,
		in.StockQuantity,
		0,
		in.Properties,
		time.Now(),
	)

	if err := s.partRepo.Create(ctx, newPart); err != nil {
		return uuid.Nil, fmt.Errorf("создать деталь: %w", err)
	}

	return partUUID, nil
}
