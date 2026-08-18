package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderapiv1 "github.com/mbakhodurov/homeworks2/week7/order/internal/api/order/v1"
	iamclient "github.com/mbakhodurov/homeworks2/week7/order/internal/client/grpc/iam/v1"
	inventoryclient "github.com/mbakhodurov/homeworks2/week7/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/mbakhodurov/homeworks2/week7/order/internal/client/grpc/payment/v1"
	"github.com/mbakhodurov/homeworks2/week7/order/internal/config"
	assemblyconsumer "github.com/mbakhodurov/homeworks2/week7/order/internal/consumer/assembly_consumer"
	"github.com/mbakhodurov/homeworks2/week7/order/internal/interceptor"
	"github.com/mbakhodurov/homeworks2/week7/order/internal/middleware"
	orderpayproducer "github.com/mbakhodurov/homeworks2/week7/order/internal/producer/order_producer"
	orderrepo "github.com/mbakhodurov/homeworks2/week7/order/internal/repository/order"
	orderitemrepo "github.com/mbakhodurov/homeworks2/week7/order/internal/repository/order_item"
	ordersvc "github.com/mbakhodurov/homeworks2/week7/order/internal/service/order"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/closer"
	kafkaconsumer "github.com/mbakhodurov/homeworks2/week7/platform/pkg/kafka/consumer"
	kafkaproducer "github.com/mbakhodurov/homeworks2/week7/platform/pkg/kafka/producer"
	kafkamw "github.com/mbakhodurov/homeworks2/week7/platform/pkg/middleware/kafka"
	orderv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/openapi/order/v1"
	authv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/payment/v1"
	"github.com/mbakhodurov/homeworks2/week7/shared/static"
)

// diContainer — контейнер зависимостей с ленивой инициализацией.
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует результат.
type diContainer struct {
	pgPool    *pgxpool.Pool
	txManager ordersvc.TxManager

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
	iamConn       *grpc.ClientConn

	orderRepo       ordersvc.OrderRepository
	orderItemRepo   ordersvc.OrderItemRepository
	inventoryClient ordersvc.InventoryClient
	paymentClient   ordersvc.PaymentClient
	iamClient       authv1.AuthServiceClient

	orderService orderapiv1.OrderService

	httpHandler http.Handler

	syncProducer      sarama.SyncProducer
	orderPaidProducer ordersvc.OrderPaidProducer
	consumerGroup     sarama.ConsumerGroup
	assemblyConsumer  *assemblyconsumer.Consumer
}

// PGPool возвращает пул подключений к PostgreSQL.
// При первом вызове создаёт пул, проверяет соединение и регистрирует closer.
func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		if err = pool.Ping(ctx); err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

// TxManager возвращает менеджер транзакций.
func (d *diContainer) TxManager(ctx context.Context) ordersvc.TxManager {
	if d.txManager == nil {
		m, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = m
	}

	return d.txManager
}

// InventoryConn возвращает gRPC-соединение с InventoryService.
// SessionForwarder прокидывает session_uuid из контекста запроса в исходящую
// metadata — без него InventoryService отклонит вызов с Unauthenticated.
func (d *diContainer) InventoryConn() *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithUnaryInterceptor(interceptor.SessionForwarder()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			os.Exit(1)
		}

		closer.Add("InventoryService gRPC-соединение", func(_ context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

// PaymentConn возвращает gRPC-соединение с PaymentService.
func (d *diContainer) PaymentConn() *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			os.Exit(1)
		}

		closer.Add("PaymentService gRPC-соединение", func(_ context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

// OrderRepository возвращает репозиторий заказов.
func (d *diContainer) OrderRepository(ctx context.Context) ordersvc.OrderRepository {
	if d.orderRepo == nil {
		d.orderRepo = orderrepo.NewRepository(d.PGPool(ctx))
	}

	return d.orderRepo
}

// OrderItemRepository возвращает репозиторий позиций заказа.
func (d *diContainer) OrderItemRepository(ctx context.Context) ordersvc.OrderItemRepository {
	if d.orderItemRepo == nil {
		d.orderItemRepo = orderitemrepo.NewRepository(d.PGPool(ctx))
	}

	return d.orderItemRepo
}

// InventoryClient возвращает клиент InventoryService.
func (d *diContainer) InventoryClient() ordersvc.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryclient.New(inventoryv1.NewInventoryServiceClient(d.InventoryConn()))
	}

	return d.inventoryClient
}

// PaymentClient возвращает клиент PaymentService.
func (d *diContainer) PaymentClient() ordersvc.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentclient.New(paymentv1.NewPaymentServiceClient(d.PaymentConn()))
	}

	return d.paymentClient
}

// IAMConn возвращает gRPC-соединение с IAMService.
func (d *diContainer) IAMConn() *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			slog.Error("не удалось подключиться к IAMService", "error", err)
			os.Exit(1)
		}

		closer.Add("IAMService gRPC-соединение", func(_ context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn
}

// IAMClient возвращает proto-клиент AuthService (IAMService).
func (d *diContainer) IAMClient() authv1.AuthServiceClient {
	if d.iamClient == nil {
		d.iamClient = authv1.NewAuthServiceClient(d.IAMConn())
	}

	return d.iamClient
}

// SyncProducer создаёт sarama.SyncProducer для отправки сообщений в Kafka.
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

// OrderPaidProducer возвращает продюсер события OrderPaid.
func (d *diContainer) OrderPaidProducer() ordersvc.OrderPaidProducer {
	if d.orderPaidProducer == nil {
		topic := config.AppConfig().OrderPaidProducer.Topic
		p := kafkaproducer.NewProducer(d.SyncProducer(), topic)
		d.orderPaidProducer = orderpayproducer.New(p)
	}
	return d.orderPaidProducer
}

// ConsumerGroup создаёт sarama.ConsumerGroup для потребления сообщений из Kafka.
func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

		brokers := config.AppConfig().Kafka.Brokers
		groupID := config.AppConfig().ShipAssembledConsumer.GroupID

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

// AssemblyConsumer возвращает consumer для ShipAssembled событий.
func (d *diContainer) AssemblyConsumer(ctx context.Context) *assemblyconsumer.Consumer {
	if d.assemblyConsumer == nil {
		topic := config.AppConfig().ShipAssembledConsumer.Topic
		topics := []string{topic}

		c := kafkaconsumer.NewConsumer(
			d.ConsumerGroup(),
			topics,
			kafkaconsumer.WithMiddlewares(kafkamw.ConsumerSession(), kafkamw.ConsumerLogging()),
		)

		d.assemblyConsumer = assemblyconsumer.NewConsumer(
			c,
			d.OrderRepository(ctx),
			d.InventoryClient(),
			d.TxManager(ctx),
		)
	}
	return d.assemblyConsumer
}

// OrderService возвращает application-сервис заказов.
func (d *diContainer) OrderService(ctx context.Context) orderapiv1.OrderService {
	if d.orderService == nil {
		d.orderService = ordersvc.New(
			d.OrderRepository(ctx),
			d.OrderItemRepository(ctx),
			d.InventoryClient(),
			d.PaymentClient(),
			d.TxManager(ctx),
			d.OrderPaidProducer(),
		)
	}

	return d.orderService
}

// HTTPHandler возвращает HTTP-обработчик OrderService (ogen-сервер поверх OpenAPI-спеки,
// обёрнутый в мидлварь аутентификации).
func (d *diContainer) HTTPHandler(ctx context.Context) http.Handler {
	if d.httpHandler == nil {
		api := orderapiv1.New(d.OrderService(ctx))

		server, err := orderv1.NewServer(api, orderv1.WithErrorHandler(orderapiv1.ErrorHandler))
		if err != nil {
			slog.Error("не удалось создать HTTP-сервер OpenAPI", "error", err)
			os.Exit(1)
		}

		iamClient := iamclient.New(d.IAMClient())

		mux := http.NewServeMux()

		mux.Handle("/api/", middleware.Auth(iamClient)(server))

		mux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, _ *http.Request) {
			data, readErr := static.FS.ReadFile("swagger-ui.html")
			if readErr != nil {
				http.Error(w, "swagger-ui.html не найден", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if _, writeErr := w.Write(data); writeErr != nil {
				slog.Error("ошибка записи swagger-ui", "error", writeErr)
			}
		})

		mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
			data, readErr := static.FS.ReadFile("generated/order.swagger.yaml")
			if readErr != nil {
				http.Error(w, "swagger.json не найден", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/yaml")
			if _, writeErr := w.Write(data); writeErr != nil {
				slog.Error("ошибка записи swagger.json", "error", writeErr)
			}
		})

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
		})

		d.httpHandler = mux
	}

	return d.httpHandler
}
