package config

// IAMClientConfig — параметры подключения к IAMService.
type IAMClientConfig struct {
	Address string `yaml:"address" env:"IAM_ADDRESS" env-default:"localhost:50053"`
}
