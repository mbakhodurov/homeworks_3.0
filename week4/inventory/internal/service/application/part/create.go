package part

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/service/input"
)

// Create создаёт новую деталь и возвращает её UUID.
// UUID и CreatedAt — идентичность записи, их назначает сервис, а не вызывающий код.
// Свойства (properties) новой детали не задаются через этот метод — он покрывает
// только базовое CRUD-создание; для типоспецифичных properties нужен отдельный конструктор.
func (s *service) Create(ctx context.Context, in input.CreatePartInput) (uuid.UUID, error) {
	partUUID := uuid.New()

	newPart := model.RestorePart(
		partUUID,
		in.Name,
		in.Description,
		in.PartType,
		in.Price,
		in.StockQuantity,
		0,
		model.PartProperties{},
		time.Now(),
	)

	if err := s.partRepo.Create(ctx, newPart); err != nil {
		return uuid.Nil, fmt.Errorf("создать деталь: %w", err)
	}

	return partUUID, nil
}
