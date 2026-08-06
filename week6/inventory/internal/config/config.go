package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultConfigPath = "config.local.yaml"

var appConfig *Config

// Config — конфигурация приложения InventoryService.
type Config struct {
	GRPC      GRPCConfig      `yaml:"grpc"`
	HTTP      HTTPConfig      `yaml:"http"`
	PG        PGConfig        `yaml:"pg"`
	IAMClient IAMClientConfig `yaml:"iam_client"`
	Logger    LoggerConfig    `yaml:"logger"`
}

// MustLoad загружает конфиг (YAML + env) и сохраняет в глобальную переменную.
// При ошибке — паникует.
func MustLoad() {
	cfgPath := resolveConfigPath()

	var cfg Config

	var err error
	if cfgPath != "" {
		err = cleanenv.ReadConfig(cfgPath, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}

	if err != nil {
		panic(fmt.Sprintf("не удалось загрузить конфиг: %v", err))
	}

	appConfig = &cfg
}

// AppConfig возвращает загруженный конфиг.
func AppConfig() *Config {
	return appConfig
}

// resolveConfigPath определяет путь к конфиг-файлу по цепочке приоритетов:
// флаг -config > env CONFIG_PATH > "config.local.yaml".
func resolveConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	return defaultConfigPath
}
