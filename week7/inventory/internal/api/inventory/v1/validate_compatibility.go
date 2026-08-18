package v1

import (
	"context"

	inventoryv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/inventory/v1"
)

// ValidateCompatibility проверяет совместимость деталей между собой.
func (a *api) ValidateCompatibility(ctx context.Context, req *inventoryv1.ValidateCompatibilityRequest) (*inventoryv1.ValidateCompatibilityResponse, error) {
	uuids := []string{req.GetHullUuid(), req.GetEngineUuid()}
	if s := req.GetShieldUuid(); s != "" {
		uuids = append(uuids, s)
	}
	if w := req.GetWeaponUuid(); w != "" {
		uuids = append(uuids, w)
	}

	if err := a.partService.ValidateCompatibility(ctx, uuids); err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
