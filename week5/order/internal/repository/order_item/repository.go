package orderitem

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// repository — PostgreSQL-хранилище позиций заказа.
type repository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

// NewRepository создаёт новый PostgreSQL-репозиторий позиций заказа.
func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}
