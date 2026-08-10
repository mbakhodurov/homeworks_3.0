package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/mbakhodurov/homeworks2/week6/inventory/internal/errors"
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
	case errors.Is(err, errs.ErrInvalidUUID),
		errors.Is(err, errs.ErrInvalidPartInfo),
		errors.Is(err, errs.ErrInvalidProperties),
		errors.Is(err, errs.ErrPartTypeMismatch):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrIncompatibleParts),
		errors.Is(err, errs.ErrNothingToRelease),
		errors.Is(err, errs.ErrNothingToCommit):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrOutOfStock):
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	default:
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}

		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
