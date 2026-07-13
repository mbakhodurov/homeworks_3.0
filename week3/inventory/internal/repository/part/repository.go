package part

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool

	// getter извлекает активную транзакцию (pgx.Tx) из context.Context
	// Подробнее о том, как работает CtxGetter и зачем нужен кастомный — см. README.md
	getter *trmpgx.CtxGetter
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}
