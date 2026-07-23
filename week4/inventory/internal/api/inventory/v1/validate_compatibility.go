package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/api/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// ValidateCompatibility проверяет совместимость деталей между собой.
func (a *api) ValidateCompatibility(ctx context.Context, req *inventoryv1.ValidateCompatibilityRequest) (*inventoryv1.ValidateCompatibilityResponse, error) {
	slots, err := converter.ToShipSlots(req)
	if err != nil {
		return nil, err
	}

	if err = a.partService.ValidateCompatibility(ctx, slots); err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
