package assembly_consumer

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/mbakhodurov/homeworks2/week7/order/internal/model"
	eventsv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/events/v1"
)

// decodeShipAssembled десериализует protobuf-сообщение в доменное событие.
func decodeShipAssembled(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("десериализовать ShipAssembled: %w", err)
	}

	eventUUID, err := uuid.Parse(pb.GetEventUuid())
	if err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("распарсить event_uuid: %w", err)
	}

	orderUUID, err := uuid.Parse(pb.GetOrderUuid())
	if err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("распарсить order_uuid: %w", err)
	}

	userUUID, err := uuid.Parse(pb.GetUserUuid())
	if err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("распарсить user_uuid: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    eventUUID,
		OrderUUID:    orderUUID,
		UserUUID:     userUUID,
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
