package config

// LoggerConfig — параметры логирования.
type LoggerConfig struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
}
