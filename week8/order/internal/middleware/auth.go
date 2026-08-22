package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week8/platform/pkg/auth"
)

const bearerPrefix = "Bearer "

// AuthClient определяет контракт, нужный HTTP-мидлвари для проверки сессии.
type AuthClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

// errorResponse — тело ответа при ошибке аутентификации, в том же формате,
// что и остальные ошибки OrderService (см. api/order/v1/error_handler.go).
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Auth — HTTP-мидлварь аутентификации. Читает Authorization: Bearer <session_uuid>
// (или заголовок X-Session-UUID), проверяет сессию через AuthService.Whoami и
// кладёт user_uuid/session_uuid в контекст запроса.
func Auth(authClient AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionUUID, ok := extractSessionUUID(r)
			if !ok {
				writeUnauthorized(w, "отсутствует заголовок Authorization")
				return
			}

			if sessionUUID == "" {
				writeUnauthorized(w, "неверный формат Authorization")
				return
			}

			userUUID, err := authClient.Whoami(r.Context(), sessionUUID)
			if err != nil {
				writeUnauthorized(w, "недействительная сессия")
				return
			}

			ctx := auth.WithUserUUID(r.Context(), userUUID)
			ctx = auth.WithSessionUUID(ctx, sessionUUID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractSessionUUID достаёт session_uuid из заголовка Authorization (Bearer)
// или, как запасной вариант, из X-Session-UUID.
//
// Возвращает (значение, найден ли вообще подходящий заголовок). Если заголовок
// Authorization присутствует, но не в формате "Bearer <uuid>" — возвращает
// ("", true), чтобы вызывающий код отличил «нет заголовка» от «неверный формат».
func extractSessionUUID(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return "", true
		}

		return strings.TrimPrefix(authHeader, bearerPrefix), true
	}

	if sessionHeader := r.Header.Get("X-Session-UUID"); sessionHeader != "" {
		return sessionHeader, true
	}

	return "", false
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(errorResponse{Code: http.StatusUnauthorized, Message: message}) //nolint:gosec // ошибка записи в ResponseWriter не восстановима
}
