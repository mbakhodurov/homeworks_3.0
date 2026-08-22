package user

import "github.com/jackc/pgx/v5/pgxpool"

// repository — PostgreSQL-хранилище пользователей.
type repository struct {
	pool *pgxpool.Pool
}

// NewRepository создаёт новый PostgreSQL-репозиторий пользователей.
func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
