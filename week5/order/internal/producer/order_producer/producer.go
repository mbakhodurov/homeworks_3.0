package order_producer

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/mbakhodurov/homeworks2/week5/order/internal/model"
	"github.com/mbakhodurov/homeworks2/week5/platform/pkg/kafka"
	kafkaproducer "github.com/mbakhodurov/homeworks2/week5/platform/pkg/kafka/producer"
	eventsv1 "github.com/mbakhodurov/homeworks2/week5/shared/pkg/proto/events/v1"
)

// Producer публикует событие OrderPaid в Kafka.
type Producer struct {
	producer *kafkaproducer.Producer
}

// New создаёт Producer для топика OrderPaid.
func New(p *kafkaproducer.Producer) *Producer {
	return &Producer{producer: p}
}

// Produce сериализует событие в protobuf и отправляет в Kafka.
func (p *Producer) Produce(ctx context.Context, event model.OrderPaidEvent) error {
	pb := &eventsv1.OrderPaid{
		EventUuid: event.EventUUID.String(),
		OrderUuid: event.OrderUUID.String(),
		UserUuid:  event.UserUUID.String(),
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("сериализовать OrderPaid: %w", err)
	}

	return p.producer.Send(ctx, &kafka.Message{
		Key:   []byte(event.OrderUUID.String()),
		Value: data,
	})
}
