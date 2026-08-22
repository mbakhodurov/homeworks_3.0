package config

// KafkaConfig — конфигурация Kafka-брокеров.
type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
}
