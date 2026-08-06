package v1

import (
	inventoryv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer

	partService PartService
}

// New создаёт новый gRPC API InventoryService.
func New(partService PartService) *api {
	return &api{
		partService: partService,
	}
}
