package config

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
}
