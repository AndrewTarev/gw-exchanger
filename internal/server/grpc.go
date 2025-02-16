package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"

	config "gw-exchanger/internal/config"
	"gw-exchanger/internal/delivery/grpc_delivery"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	exch "github.com/AndrewTarev/proto-repo/gen/exchange"

	"gw-exchanger/internal/service"
)

// RunGRPCServer запускает gRPC сервер
func RunGRPCServer(ctx context.Context, svc *service.Service, cfg *config.Config, logger *logrus.Logger) error {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.WithFields(logrus.Fields{"addr": addr, "error": err}).Error("Failed to listen")
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// Создаем gRPC сервер с middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)), // Логирование запросов
			grpcRecovery.UnaryServerInterceptor(),                       // Обработка паники
		)),
	)

	// Регистрируем сервис
	exch.RegisterExchangeServiceServer(grpcServer, grpc_delivery.NewExchangerHandler(svc))

	// Включаем gRPC Reflection (удобно для отладки)
	reflection.Register(grpcServer)

	// Канал для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		logger.WithField("addr", addr).Info("🚀 gRPC Server is running")
		if err := grpcServer.Serve(listener); err != nil {
			logger.WithField("error", err).Fatal("gRPC Server crashed")
		}
	}()

	// Блокируем выполнение, пока не придёт сигнал завершения
	select {
	case <-ctx.Done():
		logger.Warn("Context cancelled, shutting down gRPC Server...")
	case <-stop:
		logger.Info("Received shutdown signal, stopping gRPC Server...")
	}

	// Вызываем shutdown
	return ShutdownGRPCServer(grpcServer, logger)
}

// ShutdownGRPCServer плавно останавливает gRPC сервер
func ShutdownGRPCServer(grpcServer *grpc.Server, logger *logrus.Logger) error {
	// Устанавливаем таймаут для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("✅ gRPC Server stopped gracefully")
	case <-ctx.Done():
		logger.Warn("⏳ Timeout: Forcefully stopping gRPC Server")
		grpcServer.Stop() // Принудительная остановка
	}

	return nil
}
