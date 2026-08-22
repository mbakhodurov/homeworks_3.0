package config

import "fmt"

// PostgresConfig — параметры подключения к PostgreSQL.
type PostgresConfig struct {
	Host     string `yaml:"host"     env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     string `yaml:"port"     env:"POSTGRES_PORT"     env-default:"5434"`
	Database string `yaml:"database" env:"POSTGRES_DB"       env-default:"iam-service"`
	User     string `yaml:"user"     env:"POSTGRES_USER"     env-default:"iam-service-user"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"iam-service-password"`
	SSLMode  string `yaml:"sslmode"  env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}
