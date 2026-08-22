package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	errs "github.com/mbakhodurov/homeworks2/week8/inventory/internal/errors"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/inventory/v1"
)

// DeletePart удаляет деталь по UUID.
func (a *api) DeletePart(ctx context.Context, req *inventoryv1.DeletePartRequest) (*emptypb.Empty, error) {
	partUUID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, errs.ErrInvalidUUID
	}

	if err = a.partService.Delete(ctx, partUUID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
