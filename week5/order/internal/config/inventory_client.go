package config

// InventoryClientConfig — параметры подключения к InventoryService.
type InventoryClientConfig struct {
	Address string `yaml:"address" env:"INVENTORY_ADDRESS" env-default:"localhost:50051"`
}
