package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/mbakhodurov/homeworks2/week3/inventory/internal/errors"
)

// ErrorInterceptor — унарный gRPC interceptor, который маппит доменные ошибки
// сервиса в коды gRPC. Хендлеры возвращают доменную ошибку как есть,
// interceptor — единственное место, знающее про коды протокола.
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
	case errors.Is(err, errs.ErrPartNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, errs.ErrInvalidUUID), errors.Is(err, errs.ErrInvalidPartInfo):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
