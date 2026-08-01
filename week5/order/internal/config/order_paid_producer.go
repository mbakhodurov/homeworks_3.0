package config

// OrderPaidProducerConfig — конфигурация продюсера события OrderPaid.
type OrderPaidProducerConfig struct {
	Topic string `yaml:"topic" env:"KAFKA_ORDER_PAID_TOPIC" env-default:"order.paid"`
}
