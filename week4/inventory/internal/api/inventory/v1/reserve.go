package v1

import (
	"context"

	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// ReserveParts резервирует детали для заказа.
func (a *api) ReserveParts(ctx context.Context, req *inventoryv1.ReservePartsRequest) (*inventoryv1.ReservePartsResponse, error) {
	if err := a.partService.ReserveParts(ctx, req.GetUuids()); err != nil {
		return nil, err
	}

	return &inventoryv1.ReservePartsResponse{}, nil
}
