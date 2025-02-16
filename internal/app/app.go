package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	config "gw-exchanger/internal/config"
	"gw-exchanger/internal/server"
	"gw-exchanger/internal/service"
	"gw-exchanger/internal/storage"
	"gw-exchanger/pkg/db"
)

func StartApplication(cfg *config.Config, logger *logrus.Logger) error {
	// Создаем контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключение к базе данных
	dbConn, err := db.ConnectPostgres(cfg.Database.Dsn)
	if err != nil {
		logger.Fatalf("❌ Database connection failed: %v", err)
		return err
	}
	defer dbConn.Close()

	// Запускаем миграции
	// db.ApplyMigrations(cfg.Database.Dsn, cfg.Database.MigratePath)

	// Создаем зависимости
	repo := storage.NewStorage(dbConn)
	services := service.NewExchangerService(repo)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		if err := server.RunGRPCServer(ctx, services, cfg, logger); err != nil {
			logger.Fatalf("❌ gRPC server failed: %v", err)
		}
	}()

	// Перехватываем SIGTERM и SIGINT для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнал завершения
	<-stop
	logger.Warn("⚠️ Received termination signal, shutting down...")

	// Вызываем `cancel()` для graceful shutdown всех сервисов
	cancel()

	logger.Info("✅ Application shutdown complete")
	return nil
}
