package app

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authapiv1 "github.com/mbakhodurov/homeworks2/week6/iam/internal/api/auth/v1"
	userapiv1 "github.com/mbakhodurov/homeworks2/week6/iam/internal/api/user/v1"
	"github.com/mbakhodurov/homeworks2/week6/iam/internal/interceptor"
	sessionrepo "github.com/mbakhodurov/homeworks2/week6/iam/internal/repository/session"
	userrepo "github.com/mbakhodurov/homeworks2/week6/iam/internal/repository/user"
	iamsvc "github.com/mbakhodurov/homeworks2/week6/iam/internal/service/iam"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/grpc/health"
	authv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/auth/v1"
	userv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/user/v1"
)

// NewGRPCServer строит IAM-сервис из готового pgxpool.Pool и redis.Client и
// регистрирует AuthService + UserService на новом grpc.Server.
//
// Экспортирована (а не приватная деталь initGRPCServer), т.к. это единственная
// точка правды для сборки сервера — переиспользуется и здесь (App.initGRPCServer,
// зависимости берутся из конфига), и API-тестами внутри модуля iam (iam/tests),
// которым нужен IAM-сервер в bufconn с явно заданными pool/rdb.
func NewGRPCServer(pool *pgxpool.Pool, rdb *redis.Client, sessionTTL time.Duration, bcryptCost int) *grpc.Server {
	userRepo := userrepo.NewRepository(pool)
	sessionRepo := sessionrepo.NewRepository(rdb)

	svc := iamsvc.New(userRepo, sessionRepo, sessionTTL, bcryptCost)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.ErrorInterceptor),
	)

	reflection.Register(server)
	health.RegisterService(server)

	authv1.RegisterAuthServiceServer(server, authapiv1.New(svc))
	userv1.RegisterUserServiceServer(server, userapiv1.New(svc))

	return server
}
