package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week8/iam/internal/model"
	"github.com/mbakhodurov/homeworks2/week8/iam/internal/repository/redis_view"
)

// SessionToRedisView конвертирует доменную модель сессии в Redis view.
func SessionToRedisView(s model.Session) redis_view.SessionRedisView {
	return redis_view.SessionRedisView{
		UUID:      s.UUID.String(),
		UserUUID:  s.UserUUID.String(),
		Login:     s.Login,
		CreatedAt: s.CreatedAt.UnixNano(),
		ExpiresAt: s.ExpiresAt.UnixNano(),
	}
}

// SessionFromRedisView восстанавливает доменную модель сессии из Redis view.
func SessionFromRedisView(v redis_view.SessionRedisView) model.Session {
	return model.Session{
		UUID:      uuid.MustParse(v.UUID),
		UserUUID:  uuid.MustParse(v.UserUUID),
		Login:     v.Login,
		CreatedAt: time.Unix(0, v.CreatedAt),
		ExpiresAt: time.Unix(0, v.ExpiresAt),
	}
}
