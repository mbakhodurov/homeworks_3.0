package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/mbakhodurov/homeworks2/week6/iam/internal/config"
	"github.com/mbakhodurov/homeworks2/week6/platform/pkg/closer"
	platformRedis "github.com/mbakhodurov/homeworks2/week6/platform/pkg/redis"
)

// diContainer — контейнер зависимостей с ленивой инициализацией.
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует результат.
type diContainer struct {
	pgPool      *pgxpool.Pool
	redisClient *redis.Client
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
