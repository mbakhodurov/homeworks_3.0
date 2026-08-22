package v1

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	inventoryv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/inventory/v1"
)

// CommitParts списывает зарезервированные детали со склада после сборки корабля.
func (a *api) CommitParts(ctx context.Context, req *inventoryv1.CommitPartsRequest) (*emptypb.Empty, error) {
	if err := a.partService.CommitParts(ctx, req.GetUuids()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
