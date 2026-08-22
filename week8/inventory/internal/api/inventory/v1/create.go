package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/inventory/internal/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/inventory/v1"
)

// CreateParts создаёт новую деталь.
func (a *api) CreateParts(ctx context.Context, req *inventoryv1.CreatePartsRequest) (*inventoryv1.CreatePartsResponse, error) {
	in, err := converter.ToCreatePartInput(req)
	if err != nil {
		return nil, err
	}

	partUUID, err := a.partService.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.CreatePartsResponse{Uuid: partUUID.String()}, nil
}
