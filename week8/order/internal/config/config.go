package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultConfigPath = "config.local.yaml"

var appConfig *Config

// Config — конфигурация приложения OrderService.
type Config struct {
	HTTP                  HTTPConfig                  `yaml:"http"`
	PG                    PGConfig                    `yaml:"pg"`
	InventoryClient       InventoryClientConfig       `yaml:"inventory_client"`
	PaymentClient         PaymentClientConfig         `yaml:"payment_client"`
	IAMClient             IAMClientConfig             `yaml:"iam_client"`
	Logger                LoggerConfig                `yaml:"logger"`
	Kafka                 KafkaConfig                 `yaml:"kafka"`
	OrderPaidProducer     OrderPaidProducerConfig     `yaml:"order_paid_producer"`
	ShipAssembledConsumer ShipAssembledConsumerConfig `yaml:"ship_assembled_consumer"`
	OTel                  otelConfig                  `yaml:"otel"`
	RateLimit             rateLimitConfig             `yaml:"rate_limit"`
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
