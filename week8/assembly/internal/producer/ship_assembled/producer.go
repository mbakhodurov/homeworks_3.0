package ship_assembled

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mbakhodurov/homeworks2/week8/assembly/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/platform/pkg/kafka"
	kafkaproducer "github.com/mbakhodurov/homeworks2/week8/platform/pkg/kafka/producer"
	kafkamw "github.com/mbakhodurov/homeworks2/week8/platform/pkg/middleware/kafka"
	eventsv1 "github.com/mbakhodurov/homeworks2/week8/shared/pkg/proto/events/v1"
)

type Producer struct {
	producer *kafkaproducer.Producer
}

func New(p *kafkaproducer.Producer) *Producer {
	return &Producer{producer: p}
}

func (p *Producer) Produce(ctx context.Context, event model.ShipAssembledEvent) error {
	pb := &eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID.String(),
		OrderUuid:    event.OrderUUID.String(),
		UserUuid:     event.UserUUID.String(),
		BuildTimeSec: event.BuildTimeSec,
		AssembledAt:  timestamppb.New(event.AssembledAt),
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("сериализовать ShipAssembled: %w", err)
	}

	return p.producer.Send(ctx, &kafka.Message{
		Key:     []byte(event.OrderUUID.String()),
		Value:   data,
		Headers: kafkamw.ProducerSessionHeaders(ctx),
	})
}
