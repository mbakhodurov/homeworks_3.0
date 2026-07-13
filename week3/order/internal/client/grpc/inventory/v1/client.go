package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mbakhodurov/homeworks2/week3/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/mbakhodurov/homeworks2/week3/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week3/order/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week3/shared/pkg/proto/inventory/v1"
)

// client — обёртка над proto-клиентом InventoryService.
type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

// New создаёт новый клиент InventoryService.
func New(inventoryClient inventoryv1.InventoryServiceClient) *client {
	return &client{inventoryClient: inventoryClient}
}

// ListParts возвращает детали по списку UUID.
func (c *client) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	resp, err := c.inventoryClient.ListParts(ctx, converter.ToListPartsRequest(uuids))
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, errs.ErrPartNotFound
		}
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return converter.PartsToModel(resp.GetPart()), nil
}
