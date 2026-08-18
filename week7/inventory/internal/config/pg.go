package config

import "fmt"

// PGConfig — параметры подключения к PostgreSQL.
type PGConfig struct {
	Host     string `yaml:"host"     env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     string `yaml:"port"     env:"POSTGRES_PORT"     env-default:"5433"`
	Database string `yaml:"database" env:"POSTGRES_DB"       env-default:"inventory"`
	User     string `yaml:"user"     env:"POSTGRES_USER"     env-default:"inventory_admin"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"inventory_secret"`
	SSLMode  string `yaml:"sslmode"  env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *PGConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}
