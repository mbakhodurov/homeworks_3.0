package assembly_consumer

import (
	"context"

	kafkaconsumer "github.com/mbakhodurov/homeworks2/week8/platform/pkg/kafka/consumer"
)

// Consumer запускает обработку событий ShipAssembled из Kafka.
type Consumer struct {
	consumer *kafkaconsumer.Consumer
	handler  *handler
}

// NewConsumer создаёт Consumer для обработки ShipAssembled событий.
func NewConsumer(
	consumer *kafkaconsumer.Consumer,
	orderRepo OrderRepository,
	inventoryClient InventoryClient,
	txManager TxManager,
) *Consumer {
	return &Consumer{
		consumer: consumer,
		handler:  newHandler(orderRepo, inventoryClient, txManager),
	}
}

// Run запускает бесконечный цикл потребления. Блокирует до отмены ctx.
func (c *Consumer) Run(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.handler.Handle)
}
