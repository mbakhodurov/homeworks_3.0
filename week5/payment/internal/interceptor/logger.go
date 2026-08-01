package interceptor

import (
	"context"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
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
