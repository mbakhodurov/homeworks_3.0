package ratelimit

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

// Middleware создаёт HTTP middleware с распределённым rate limiter.
//
// Библиотека redis_rate использует алгоритм GCRA (Generic Cell Rate Algorithm) —
// вариант token bucket, реализованный атомарно в Redis через Lua-скрипт.
// Благодаря Redis все инстансы сервиса разделяют общий счётчик запросов,
// поэтому лимит работает глобально, а не per-instance.
//
// Ключ для лимита — путь запроса (r.URL.Path): разные эндпоинты получают
// независимые счётчики.
//
// Стратегия fail-open: если Redis недоступен — запрос пропускается.
// Лучше временно пропустить лишний трафик, чем остановить весь сервис
// из-за проблем с Redis.
func Middleware(limiter *redis_rate.Limiter, rate, burst int) func(http.Handler) http.Handler {
	limit := redis_rate.Limit{
		Rate:   rate,
		Burst:  burst,
		Period: time.Second,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Path

			res, err := limiter.Allow(r.Context(), key, limit)
			if err != nil {
				slog.WarnContext(r.Context(), "ошибка проверки rate limit, пропускаем запрос",
					"path", key,
					"error", err,
				)
				next.ServeHTTP(w, r)
				return
			}

			if res.Allowed == 0 {
				slog.WarnContext(r.Context(), "запрос отклонён rate limiter",
					"path", key,
					"retry_after", res.RetryAfter,
				)
				http.Error(w, "слишком много запросов, попробуйте позже", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
