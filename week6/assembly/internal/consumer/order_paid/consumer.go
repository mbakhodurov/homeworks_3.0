package order_paid

import (
	"context"

	kafkaconsumer "github.com/mbakhodurov/homeworks2/week6/platform/pkg/kafka/consumer"
)

type Consumer struct {
	consumer *kafkaconsumer.Consumer
	handler  *handler
}

func NewConsumer(consumer *kafkaconsumer.Consumer, assemblyService AssemblyService) *Consumer {
	return &Consumer{
		consumer: consumer,
		handler:  newHandler(assemblyService),
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.handler.Handle)
}
