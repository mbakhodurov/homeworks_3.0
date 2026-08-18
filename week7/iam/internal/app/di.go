package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	authapiv1 "github.com/mbakhodurov/homeworks2/week7/iam/internal/api/auth/v1"
	userapiv1 "github.com/mbakhodurov/homeworks2/week7/iam/internal/api/user/v1"
	"github.com/mbakhodurov/homeworks2/week7/iam/internal/config"
	sessionrepo "github.com/mbakhodurov/homeworks2/week7/iam/internal/repository/session"
	userrepo "github.com/mbakhodurov/homeworks2/week7/iam/internal/repository/user"
	iamsvc "github.com/mbakhodurov/homeworks2/week7/iam/internal/service/iam"
	"github.com/mbakhodurov/homeworks2/week7/platform/pkg/closer"
	platformRedis "github.com/mbakhodurov/homeworks2/week7/platform/pkg/redis"
	authv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/auth/v1"
	userv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/user/v1"
)

// iamService объединяет контракты AuthService и UserService — оба API-обработчика
// (AuthAPI и UserAPI) работают с одним и тем же сервисом IAM.
type iamService interface {
	authapiv1.AuthService
	userapiv1.UserService
}

// diContainer — контейнер зависимостей с ленивой инициализацией.
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует результат.
type diContainer struct {
	pgPool      *pgxpool.Pool
	redisClient *redis.Client

	userRepo    iamsvc.UserRepository
	sessionRepo iamsvc.SessionRepository

	iamService iamService

	authHandler authv1.AuthServiceServer
	userHandler userv1.UserServiceServer
}

// PGPool возвращает пул подключений к PostgreSQL.
func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		if err = pool.Ping(ctx); err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

// RedisClient возвращает клиент Redis (хранилище сессий).
func (d *diContainer) RedisClient(_ context.Context) *redis.Client {
	if d.redisClient == nil {
		rdb, err := platformRedis.NewClient(&redis.Options{
			Addr:        config.AppConfig().Redis.Address(),
			Password:    config.AppConfig().Redis.Password,
			DB:          config.AppConfig().Redis.DB,
			DialTimeout: config.AppConfig().Redis.ConnectionTimeout,
		}, slog.Default())
		if err != nil {
			slog.Error("не удалось создать Redis клиент", "error", err)
			os.Exit(1)
		}

		closer.Add("Redis", func(_ context.Context) error {
			return rdb.Close()
		})

		d.redisClient = rdb
	}

	return d.redisClient
}

// UserRepository возвращает PostgreSQL-репозиторий пользователей.
func (d *diContainer) UserRepository(ctx context.Context) iamsvc.UserRepository {
	if d.userRepo == nil {
		d.userRepo = userrepo.NewRepository(d.PGPool(ctx))
	}

	return d.userRepo
}

// SessionRepository возвращает Redis-репозиторий сессий.
func (d *diContainer) SessionRepository(ctx context.Context) iamsvc.SessionRepository {
	if d.sessionRepo == nil {
		d.sessionRepo = sessionrepo.NewRepository(d.RedisClient(ctx))
	}

	return d.sessionRepo
}

// IAMService возвращает сервис IAM (аутентификация + управление пользователями).
func (d *diContainer) IAMService(ctx context.Context) iamService {
	if d.iamService == nil {
		d.iamService = iamsvc.New(
			d.UserRepository(ctx),
			d.SessionRepository(ctx),
			config.AppConfig().Session.TTL,
			bcrypt.DefaultCost,
		)
	}

	return d.iamService
}

// AuthAPI возвращает gRPC-обработчик AuthService.
func (d *diContainer) AuthAPI(ctx context.Context) authv1.AuthServiceServer {
	if d.authHandler == nil {
		d.authHandler = authapiv1.New(d.IAMService(ctx))
	}

	return d.authHandler
}

// UserAPI возвращает gRPC-обработчик UserService.
func (d *diContainer) UserAPI(ctx context.Context) userv1.UserServiceServer {
	if d.userHandler == nil {
		d.userHandler = userapiv1.New(d.IAMService(ctx))
	}

	return d.userHandler
}
