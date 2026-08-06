package config

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
}
