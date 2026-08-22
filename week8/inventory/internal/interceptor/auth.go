package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/mbakhodurov/homeworks2/week8/platform/pkg/auth"
)

// SessionMetadataKey — имя gRPC-metadata ключа, через который клиент передаёт
// session_uuid. Совпадает с ключом, который пишет order/internal/interceptor.SessionForwarder.
const SessionMetadataKey = "session-uuid"

// AuthClient определяет контракт, нужный interceptor'у для проверки сессии.
type AuthClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

// Auth — gRPC server/incoming interceptor. Читает session-uuid из incoming
// metadata, проверяет сессию через AuthService.Whoami и кладёт user_uuid
// в контекст перед вызовом хендлера.
func Auth(authClient AuthClient) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "отсутствует metadata")
		}

		values := md.Get(SessionMetadataKey)
		if len(values) == 0 || values[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "отсутствует session-uuid")
		}

		sessionUUID := values[0]

		userUUID, err := authClient.Whoami(ctx, sessionUUID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "недействительная сессия")
		}

		ctx = auth.WithUserUUID(ctx, userUUID)

		return handler(ctx, req)
	}
}
