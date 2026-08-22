package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	errs "github.com/mbakhodurov/homeworks2/week8/iam/internal/errors"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/repository/converter"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/repository/redis_view"
)

const cacheKeyPrefix = "session:"

// repository — Redis-хранилище сессий пользователей.
// Сессии хранятся как HashMap (HSET/HGETALL/DEL) — это первичное хранилище
// (источник правды для авторизации), а не кэш.
type repository struct {
	client *redis.Client
}

// NewRepository создаёт новый Redis-репозиторий сессий.
func NewRepository(client *redis.Client) *repository {
	return &repository{client: client}
}

func (r *repository) getCacheKey(sessionUUID string) string {
	return cacheKeyPrefix + sessionUUID
}

// Create сохраняет сессию в Redis (HSET) и устанавливает TTL (EXPIRE).
func (r *repository) Create(ctx context.Context, session model.Session, ttl time.Duration) error {
	cacheKey := r.getCacheKey(session.UUID.String())

	err := r.client.HSet(ctx, cacheKey, converter.SessionToRedisView(session)).Err()
	if err != nil {
		return fmt.Errorf("сохранить сессию: %w", err)
	}

	return r.client.Expire(ctx, cacheKey, ttl).Err()
}

// Get возвращает сессию по UUID (HGETALL). sessionUUID — произвольная строка:
// невалидный формат (не UUID) просто не совпадёт ни с одним ключом и даст ErrSessionNotFound.
func (r *repository) Get(ctx context.Context, sessionUUID string) (model.Session, error) {
	cacheKey := r.getCacheKey(sessionUUID)

	var view redis_view.SessionRedisView

	err := r.client.HGetAll(ctx, cacheKey).Scan(&view)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Session{}, errs.ErrSessionNotFound
		}

		return model.Session{}, fmt.Errorf("получить сессию: %w", err)
	}

	// HGetAll для несуществующего ключа возвращает пустую map без ошибки —
	// отличаем «нет такого ключа» по пустому обязательному полю UUID.
	if view.UUID == "" {
		return model.Session{}, errs.ErrSessionNotFound
	}

	return converter.SessionFromRedisView(view), nil
}

// Delete удаляет сессию (DEL). Идемпотентен — удаление несуществующей сессии не ошибка.
func (r *repository) Delete(ctx context.Context, sessionUUID string) error {
	cacheKey := r.getCacheKey(sessionUUID)

	if err := r.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("удалить сессию: %w", err)
	}

	return nil
}
