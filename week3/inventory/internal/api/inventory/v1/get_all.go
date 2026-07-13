package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/api/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/inventory/v1"
)

// GetAllPart возвращает все детали без фильтрации.
func (a *api) GetAllPart(ctx context.Context, _ *inventoryv1.GetAllPartRequest) (*inventoryv1.GetAllPartResponse, error) {
	parts, err := a.partService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dto := converter.PartsToDTO(parts)

	return &inventoryv1.GetAllPartResponse{Part: dto, TotalCount: int64(len(dto))}, nil
}
