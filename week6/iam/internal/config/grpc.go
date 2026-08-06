package config

// GRPCConfig — параметры прослушивания gRPC-сервера.
type GRPCConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"50053"`
}

// Address возвращает адрес для прослушивания gRPC-сервера.
func (c *GRPCConfig) Address() string {
	return c.Host + ":" + c.Port
}
