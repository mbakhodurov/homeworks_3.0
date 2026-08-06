package assembly

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week6/assembly/internal/model"
)

type ShipAssembledProducer interface {
	Produce(ctx context.Context, event model.ShipAssembledEvent) error
}
