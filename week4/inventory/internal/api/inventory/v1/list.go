package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей с опциональной фильтрацией.
func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	filter, err := converter.ToListFilter(req)
	if err != nil {
		return nil, err
	}

	parts, err := a.partService.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	dto := converter.PartsToDTO(parts)

	return &inventoryv1.ListPartsResponse{Part: dto, TotalCount: int64(len(dto))}, nil
}
