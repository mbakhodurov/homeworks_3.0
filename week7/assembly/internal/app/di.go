package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/IBM/sarama"

	"github.com/mbakhodurov/homeworks2/week7/assembly/internal/config"
	orderpaidconsumer "github.com/mbakhodurov/homeworks2/week7/assembly/internal/consumer/order_paid"
	shipproducer "github.com/mbakhodurov/homeworks2/week7/assembly/internal/producer/ship_assembled"
	assemblysvc "github.com/mbakhodurov/homeworks2/week7/assembly/internal/service/assembly"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/closer"
	kafkaconsumer "github.com/mbakhodurov/homeworks2/week7/platform/pkg/kafka/consumer"
	kafkaproducer "github.com/mbakhodurov/homeworks2/week7/platform/pkg/kafka/producer"
	kafkamw "github.com/mbakhodurov/homeworks2/week7/platform/pkg/middleware/kafka"
)

type diContainer struct {
	syncProducer          sarama.SyncProducer
	consumerGroup         sarama.ConsumerGroup
	shipAssembledProducer assemblysvc.ShipAssembledProducer
	assemblyService       orderpaidconsumer.AssemblyService
	orderPaidConsumer     *orderpaidconsumer.Consumer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		cfg.Producer.RequiredAcks = sarama.WaitForAll

		p, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers, cfg)
		if err != nil {
			slog.Error("не удалось создать Kafka producer", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka producer", func(_ context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
			sarama.NewBalanceStrategyRoundRobin(),
		}

		brokers := config.AppConfig().Kafka.Brokers
		groupID := config.AppConfig().OrderPaidConsumer.GroupID

		cg, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
		if err != nil {
			slog.Error("не удалось создать Kafka consumer group", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return cg.Close()
		})

		d.consumerGroup = cg
	}

	return d.consumerGroup
}

func (d *diContainer) ShipAssembledProducer() assemblysvc.ShipAssembledProducer {
	if d.shipAssembledProducer == nil {
		topic := config.AppConfig().ShipAssembledProducer.Topic
		kp := kafkaproducer.NewProducer(d.SyncProducer(), topic)
		d.shipAssembledProducer = shipproducer.New(kp)
	}

	return d.shipAssembledProducer
}

func (d *diContainer) AssemblyService() orderpaidconsumer.AssemblyService {
	if d.assemblyService == nil {
		assemblerCfg := config.AppConfig().Assembler
		d.assemblyService = assemblysvc.New(
			d.ShipAssembledProducer(),
			assemblerCfg.MinBuildTimeSec,
			assemblerCfg.MaxBuildTimeSec,
		)
	}

	return d.assemblyService
}

// OrderPaidConsumer возвращает consumer для OrderPaid событий.
// ConsumerSession вытаскивает session-uuid из Kafka-заголовка входящего сообщения
// и кладёт его в context — без него исходящий ShipAssembled не унесёт session-uuid
// дальше (см. ProducerSessionHeaders в producer/ship_assembled/producer.go).
func (d *diContainer) OrderPaidConsumer() *orderpaidconsumer.Consumer {
	if d.orderPaidConsumer == nil {
		topic := config.AppConfig().OrderPaidConsumer.Topic

		c := kafkaconsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{topic},
			kafkaconsumer.WithMiddlewares(kafkamw.ConsumerSession(), kafkamw.ConsumerLogging()),
		)

		d.orderPaidConsumer = orderpaidconsumer.NewConsumer(c, d.AssemblyService())
	}

	return d.orderPaidConsumer
}
