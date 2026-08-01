package config

type OrderPaidConsumerConfig struct {
	Topic   string `yaml:"topic"    env:"KAFKA_ORDER_PAID_TOPIC"    env-default:"order.paid"`
	GroupID string `yaml:"group_id" env:"KAFKA_ORDER_PAID_GROUP_ID" env-default:"assembly-service"`
}
