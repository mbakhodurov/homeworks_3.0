package assembly

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week7/assembly/internal/model"
)

type ShipAssembledProducer interface {
	Produce(ctx context.Context, event model.ShipAssembledEvent) error
}
