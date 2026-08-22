package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/mbakhodurov/homeworks2/week8/platform/pkg/auth"
)

// SessionMetadataKey — имя gRPC-metadata ключа, через который session_uuid
// передаётся в исходящих вызовах. Совпадает с ключом, который читает
// inventory/internal/interceptor.Auth на стороне сервера.
const SessionMetadataKey = "session-uuid"

// SessionForwarder — gRPC client/outgoing interceptor. Достаёт session_uuid
// из контекста (положен туда HTTP-мидлварью или Kafka ConsumerSession) и
// прокидывает его в исходящую gRPC-metadata. Без него InventoryService не
// получит токен сессии и отклонит вызов с Unauthenticated.
func SessionForwarder() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if sessionUUID, ok := auth.SessionUUIDFromContext(ctx); ok && sessionUUID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, SessionMetadataKey, sessionUUID)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
