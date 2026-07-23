package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/api/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
)

// ReserveParts резервирует детали для заказа.
func (a *api) ReserveParts(ctx context.Context, req *inventoryv1.ReservePartsRequest) (*inventoryv1.ReservePartsResponse, error) {
	uuids, err := converter.ToUUIDs(req.GetUuids())
	if err != nil {
		return nil, err
	}

	if err = a.partService.ReserveParts(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.ReservePartsResponse{}, nil
}
