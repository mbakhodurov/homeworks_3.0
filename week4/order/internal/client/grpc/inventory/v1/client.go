package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mbakhodurov/homeworks2/week4/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/mbakhodurov/homeworks2/week4/order/internal/errors"
	"github.com/mbakhodurov/homeworks2/week4/order/internal/model"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week4/shared/pkg/proto/inventory/v1"
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
		return nil, fmt.Errorf("получить список деталей: %w", mapGRPCError(err))
	}

	return converter.PartsToModel(resp.GetPart()), nil
}

// ValidateCompatibility проверяет совместимость деталей по слотам корабля.
func (c *client) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	_, err := c.inventoryClient.ValidateCompatibility(ctx, converter.ShipSlotsToProto(slots))
	if err != nil {
		return fmt.Errorf("проверить совместимость деталей: %w", mapGRPCError(err))
	}

	return nil
}

// ReserveParts резервирует детали для заказа.
func (c *client) ReserveParts(ctx context.Context, uuids []uuid.UUID) error {
	_, err := c.inventoryClient.ReserveParts(ctx, &inventoryv1.ReservePartsRequest{
		Uuids: converter.ToUUIDStrings(uuids),
	})
	if err != nil {
		return fmt.Errorf("зарезервировать детали: %w", mapGRPCError(err))
	}

	return nil
}

// ReleaseParts освобождает ранее зарезервированные детали.
func (c *client) ReleaseParts(ctx context.Context, uuids []uuid.UUID) error {
	_, err := c.inventoryClient.ReleaseParts(ctx, &inventoryv1.ReleasePartsRequest{
		Uuids: converter.ToUUIDStrings(uuids),
	})
	if err != nil {
		return fmt.Errorf("освободить детали: %w", mapGRPCError(err))
	}

	return nil
}

// mapGRPCError маппит gRPC-код ошибки InventoryService в доменную ошибку order.
// Коды, не входящие в таблицу, возвращаются без изменений — вызывающий код завернёт их в fmt.Errorf.
func mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return errs.ErrPartNotFound
	case codes.InvalidArgument:
		return errs.ErrPartTypeMismatch
	case codes.FailedPrecondition:
		return errs.ErrIncompatibleParts
	case codes.ResourceExhausted:
		return errs.ErrOutOfStock
	default:
		return err
	}
}
