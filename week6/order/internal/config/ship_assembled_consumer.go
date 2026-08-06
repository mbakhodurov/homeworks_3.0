package config

// ShipAssembledConsumerConfig — конфигурация консьюмера события ShipAssembled.
type ShipAssembledConsumerConfig struct {
	Topic   string `yaml:"topic"    env:"KAFKA_SHIP_ASSEMBLED_TOPIC"    env-default:"assembly.ship-assembled"`
	GroupID string `yaml:"group_id" env:"KAFKA_SHIP_ASSEMBLED_GROUP_ID" env-default:"order-service"`
}
