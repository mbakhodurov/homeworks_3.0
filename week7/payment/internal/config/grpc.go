package config

// GRPCConfig — настройки gRPC-сервера PaymentService.
type GRPCConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"50052"`
}

// Address возвращает адрес для net.Listen.
func (c *GRPCConfig) Address() string {
	return c.Host + ":" + c.Port
}
