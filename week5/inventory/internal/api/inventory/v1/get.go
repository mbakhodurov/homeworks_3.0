package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week5/inventory/internal/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/proto/inventory/v1"
)

// GetPart возвращает деталь по UUID.
func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	p, err := a.partService.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{Part: converter.PartToDTO(p)}, nil
}
