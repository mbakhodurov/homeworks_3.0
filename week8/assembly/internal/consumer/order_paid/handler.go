package order_paid

import (
	"context"
	"log/slog"

	"github.com/mbakhodurov/homeworks2/week8/platform/pkg/kafka"
)

type handler struct {
	assemblyService AssemblyService
}

func newHandler(assemblyService AssemblyService) *handler {
	return &handler{assemblyService: assemblyService}
}

func (h *handler) Handle(ctx context.Context, msg kafka.Message) error {
	event, err := decodeOrderPaid(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось десериализовать OrderPaid — пропускаем", "error", err)
		return nil
	}

	slog.InfoContext(ctx, "получено событие OrderPaid", "order_uuid", event.OrderUUID)

	return h.assemblyService.Assemble(ctx, event)
}
