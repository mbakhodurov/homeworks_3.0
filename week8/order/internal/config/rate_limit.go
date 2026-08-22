package config

// rateLimitConfig — параметры распределённого rate limiter'а.
// Приватная структура: конфигурация торчит наружу только через поле Config.RateLimit.
type rateLimitConfig struct {
	RedisAddress string `yaml:"redis_address" env:"RATE_LIMIT_REDIS_ADDRESS" env-default:"localhost:6379"`
	Rate         int    `yaml:"rate"          env:"RATE_LIMIT_RATE"          env-default:"100"`
	Burst        int    `yaml:"burst"         env:"RATE_LIMIT_BURST"         env-default:"200"`
}
