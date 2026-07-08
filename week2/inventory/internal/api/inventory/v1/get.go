package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/api/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week2/shared/pkg/proto/inventory/v1"
)

// GetPart возвращает деталь по UUID.
func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	partUUID, err := converter.ToUUID(req.GetUuid())
	if err != nil {
		return nil, err
	}

	part, err := a.partService.Get(ctx, partUUID)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{Part: converter.PartToDTO(part)}, nil
}
