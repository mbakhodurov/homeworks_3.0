package order_paid

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/mbakhodurov/homeworks2/week6/assembly/internal/model"
	eventsv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/events/v1"
)

func decodeOrderPaid(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("десериализовать OrderPaid: %w", err)
	}

	eventUUID, err := uuid.Parse(pb.GetEventUuid())
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("распарсить event_uuid: %w", err)
	}

	orderUUID, err := uuid.Parse(pb.GetOrderUuid())
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("распарсить order_uuid: %w", err)
	}

	userUUID, err := uuid.Parse(pb.GetUserUuid())
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("распарсить user_uuid: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID: eventUUID,
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
	}, nil
}
