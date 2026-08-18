package redis_view

// SessionRedisView — представление сессии под конкретное хранилище (Redis HashMap),
// а не строка таблицы — отсюда имя redis_view вместо record.
// Скалярные поля ложатся прямо в HSET без JSON-сериализации.
type SessionRedisView struct {
	UUID      string `redis:"uuid"`
	UserUUID  string `redis:"user_uuid"`
	Login     string `redis:"login"`
	CreatedAt int64  `redis:"created_at"`
	ExpiresAt int64  `redis:"expires_at"`
}
