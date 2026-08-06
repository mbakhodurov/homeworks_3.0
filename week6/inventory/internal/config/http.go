package config

// HTTPConfig — параметры HTTP-сервера (grpc-gateway + Swagger UI).
type HTTPConfig struct {
	Host string `yaml:"host" env:"HTTP_HOST" env-default:""`
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8083"`
}

// Address возвращает адрес для net.Listen.
func (c *HTTPConfig) Address() string {
	return c.Host + ":" + c.Port
}
