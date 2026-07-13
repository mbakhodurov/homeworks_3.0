package v1

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/api/converter"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/inventory/v1"
)

// DeletePart удаляет деталь по UUID.
func (a *api) DeletePart(ctx context.Context, req *inventoryv1.DeletePartRequest) (*emptypb.Empty, error) {
	partUUID, err := converter.ToUUID(req.GetUuid())
	if err != nil {
		return nil, err
	}

	if err = a.partService.Delete(ctx, partUUID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
