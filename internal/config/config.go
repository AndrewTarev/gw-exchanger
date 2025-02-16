package configs

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Конфигурация grpc сервера
type Grpc struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Конфигурация логирования
type LoggerConfig struct {
	Level       string   `mapstructure:"level"`
	Format      string   `mapstructure:"format"`
	OutputFile  string   `mapstructure:"output_file"`
	KafkaTopic  string   `mapstructure:"kafka_topic"`
	KafkaBroker []string `mapstructure:"kafka_broker"`
}

// Конфигурация базы данных
type PostgresConfig struct {
	Dsn         string `mapstructure:"dsn"`
	MigratePath string `mapstructure:"migrate_path"`
}

// Полная конфигурация
type Config struct {
	Grpc     Grpc           `mapstructure:"grpc"`
	Logging  LoggerConfig   `mapstructure:"logging"`
	Database PostgresConfig `mapstructure:"database"`
}

// LoadConfig загружает конфигурацию из файлов и переменных окружения
func LoadConfig(path string) (*Config, error) {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Инициализация Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Чтение конфигурации
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not load YAML config file: %v", err)
	}

	// Маппинг данных в структуру Config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
