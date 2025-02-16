package main

import (
	"fmt"

	"gw-exchanger/internal/app"
	configs "gw-exchanger/internal/config"
	"gw-exchanger/pkg/logging"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig("./internal/config")
	if err != nil {
		fmt.Printf("error loading config: %v", err)
	}
	// Настройка логгера
	logger, err := logging.SetupLogger(
		cfg.Logging.Level,
		cfg.Logging.Format,
		cfg.Logging.OutputFile,
		cfg.Logging.KafkaTopic,
		cfg.Logging.KafkaBroker,
	)
	if err != nil {
		logger.Fatalf("Error setting up logger: %v", err)
	}

	err = app.StartApplication(cfg, logger)
	if err != nil {
		logger.Fatalf("Error starting application: %v", err)
		return
	}
}
