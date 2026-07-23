package config

type GRPCConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"50052"`
}

func (c *GRPCConfig) Address() string {
	return c.Host + ":" + c.Port
}
