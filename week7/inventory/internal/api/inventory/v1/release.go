package v1

import (
	"context"

	inventoryv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/inventory/v1"
)

// ReleaseParts освобождает ранее зарезервированные детали.
func (a *api) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	if err := a.partService.ReleaseParts(ctx, req.GetUuids()); err != nil {
		return nil, err
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
