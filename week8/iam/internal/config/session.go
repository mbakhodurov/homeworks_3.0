package config

import "time"

// SessionConfig — параметры сессий пользователей.
// SESSION_TTL — это TTL самой сессии как источника правды для авторизации
// (не кэш, а первичное хранилище). Потеря сессии = принудительный logout пользователя.
type SessionConfig struct {
	TTL time.Duration `yaml:"ttl" env:"SESSION_TTL" env-default:"24h"`
}
