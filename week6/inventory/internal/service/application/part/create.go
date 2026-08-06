package part

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week6/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week6/inventory/internal/model"
)

// Create создаёт новую деталь и возвращает её UUID.
func (s *service) Create(ctx context.Context, in model.CreatePartInput) (uuid.UUID, error) {
	if err := validatePartTypeMatchesProperties(in.PartType, in.Properties); err != nil {
		return uuid.Nil, err
	}

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

func validatePartTypeMatchesProperties(partType model.PartType, props model.PartProperties) error {
	switch partType {
	case model.PartTypeHull:
		if props.Hull() == nil {
			return fmt.Errorf("для типа HULL ожидаются hull-свойства: %w", errs.ErrInvalidPartInfo)
		}
	case model.PartTypeEngine:
		if props.Engine() == nil {
			return fmt.Errorf("для типа ENGINE ожидаются engine-свойства: %w", errs.ErrInvalidPartInfo)
		}
	case model.PartTypeShield:
		if props.Shield() == nil {
			return fmt.Errorf("для типа SHIELD ожидаются shield-свойства: %w", errs.ErrInvalidPartInfo)
		}
	case model.PartTypeWeapon:
		if props.Weapon() == nil {
			return fmt.Errorf("для типа WEAPON ожидаются weapon-свойства: %w", errs.ErrInvalidPartInfo)
		}
	}
	return nil
}
