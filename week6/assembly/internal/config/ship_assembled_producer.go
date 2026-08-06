package config

type ShipAssembledProducerConfig struct {
	Topic string `yaml:"topic" env:"KAFKA_SHIP_ASSEMBLED_TOPIC" env-default:"assembly.ship-assembled"`
}
