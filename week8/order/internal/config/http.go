package config

// HTTPConfig — параметры прослушивания HTTP-сервера.
type HTTPConfig struct {
	Host string `yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

// Address возвращает адрес для прослушивания HTTP-сервера.
func (c *HTTPConfig) Address() string {
	return c.Host + ":" + c.Port
}
