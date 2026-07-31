package interceptor

import (
	"context"
	"errors"
	"log/slog"
	"path"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/mbakhodurov/homeworks2/week5/payment/internal/errors"
)

// RecoveryInterceptor перехватывает панику в обработчике и превращает её в codes.Internal,
// чтобы падение одного запроса не роняло весь gRPC-сервер.
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				method := path.Base(info.FullMethod)
				slog.ErrorContext(ctx, "паника в gRPC методе",
					"method", method,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				err = status.Errorf(codes.Internal, "внутренняя ошибка сервера")
			}
		}()
		return handler(ctx, req)
	}
}

// ErrorInterceptor превращает доменные ошибки в gRPC-коды.
func ErrorInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	switch {
	case errors.Is(err, errs.ErrInvalidOrderUUID), errors.Is(err, errs.ErrInvalidPaymentMethod):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}

// LoggerInterceptor логирует начало, завершение и длительность обработки каждого gRPC-запроса.
func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		method := path.Base(info.FullMethod)
		slog.InfoContext(ctx, "начало gRPC метода", "method", method)

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			slog.ErrorContext(ctx, "gRPC метод завершён с ошибкой",
				"method", method, "code", st.Code(), "error", err, "duration", duration)
		} else {
			slog.InfoContext(ctx, "gRPC метод завершён успешно", "method", method, "duration", duration)
		}

		return resp, err
	}
}
