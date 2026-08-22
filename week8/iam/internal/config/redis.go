package config

import (
	"net"
	"time"
)

// RedisConfig — параметры подключения к Redis (хранилище сессий).
type RedisConfig struct {
	Host              string        `yaml:"host"               env:"REDIS_HOST"               env-default:"localhost"`
	Port              string        `yaml:"port"               env:"REDIS_PORT"               env-default:"6379"`
	Password          string        `yaml:"password"           env:"REDIS_PASSWORD"           env-default:""`
	DB                int           `yaml:"db"                 env:"REDIS_DB"                 env-default:"0"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" env:"REDIS_CONNECTION_TIMEOUT" env-default:"5s"`
}

// Address возвращает адрес Redis в формате host:port.
func (c *RedisConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
