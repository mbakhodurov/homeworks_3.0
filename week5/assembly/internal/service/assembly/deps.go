package assembly

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week5/assembly/internal/model"
)

type ShipAssembledProducer interface {
	Produce(ctx context.Context, event model.ShipAssembledEvent) error
}
