package config

// PaymentClientConfig — параметры подключения к PaymentService.
type PaymentClientConfig struct {
	Address string `yaml:"address" env:"PAYMENT_ADDRESS" env-default:"localhost:50052"`
}
