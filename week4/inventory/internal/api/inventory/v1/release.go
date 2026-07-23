package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/api/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// ReleaseParts освобождает ранее зарезервированные детали.
func (a *api) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	uuids, err := converter.ToUUIDs(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err = a.partService.ReleaseParts(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
