package assembly

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week5/assembly/internal/model"
)

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	spread := s.maxBuildTime - s.minBuildTime
	buildTime := s.minBuildTime + time.Duration(rand.Int64N(int64(spread)))

	slog.InfoContext(ctx, "начало сборки корабля",
		"order_uuid", event.OrderUUID,
		"build_time", buildTime,
	)

	select {
	case <-time.After(buildTime):
	case <-ctx.Done():
		return ctx.Err()
	}

	assembledAt := time.Now()

	assembled := model.ShipAssembledEvent{
		EventUUID:    uuid.New(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(buildTime.Seconds()),
		AssembledAt:  assembledAt,
	}

	if err := s.producer.Produce(ctx, assembled); err != nil {
		return fmt.Errorf("отправить ShipAssembled: %w", err)
	}

	slog.InfoContext(ctx, "корабль собран",
		"order_uuid", event.OrderUUID,
		"build_time_sec", assembled.BuildTimeSec,
	)

	return nil
}
