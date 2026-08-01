package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week5/inventory/internal/converter"
	"github.com/mbakhodurov/homeworks2/week5/inventory/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/proto/inventory/v1"
)

// GetAllPart возвращает все детали без фильтрации.
func (a *api) GetAllPart(ctx context.Context, _ *inventoryv1.GetAllPartRequest) (*inventoryv1.GetAllPartResponse, error) {
	parts, err := a.partService.List(ctx, model.PartFilter{})
	if err != nil {
		return nil, err
	}

	dto := converter.PartsToDTO(parts)

	return &inventoryv1.GetAllPartResponse{Part: dto, TotalCount: int64(len(dto))}, nil
}
