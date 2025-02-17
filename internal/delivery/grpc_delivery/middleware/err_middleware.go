package middleware

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gw-exchanger/internal/errs"
)

// UnaryErrorInterceptor - middleware для обработки ошибок
func UnaryErrorInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			// Список известных ошибок, которые можно отправить клиенту
			knownErrors := map[error]codes.Code{
				errs.ErrNoRows:               codes.NotFound,
				errs.ErrUnsupportedInputCurr: codes.NotFound,
				errs.ErrUnsupportedOutputCur: codes.NotFound,
			}

			// Если ошибка есть в списке известных, отправляем клиенту как есть
			for knownErr, grpcCode := range knownErrors {
				if errors.Is(err, knownErr) {
					return nil, status.Errorf(grpcCode, err.Error())
				}
			}

			// Логируем только неизвестные ошибки
			logger.WithFields(logrus.Fields{
				"method":      info.FullMethod,
				"error":       fmt.Sprintf("%v", err),
				"stack_trace": string(debug.Stack()), // Получаем stack trace
			}).Error("❌ Unhandled gRPC error")
			return nil, status.Errorf(codes.Internal, "internal server error")
		}

		return resp, nil
	}
}

// RecoveryInterceptor - middleware для обработки паники
func RecoveryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFields(logrus.Fields{
					"method":      info.FullMethod,
					"panic":       r,
					"stack_trace": string(debug.Stack()),
				}).Error("🔥 Panic recovered in gRPC handler")

				err = status.Errorf(codes.Internal, "unexpected server error")
			}
		}()
		return handler(ctx, req)
	}
}
