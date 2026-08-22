package order_paid

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week8/assembly/internal/model"
)

type AssemblyService interface {
	Assemble(ctx context.Context, event model.OrderPaidEvent) error
}
