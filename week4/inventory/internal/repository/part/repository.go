package part

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const partColumns = `uuid, name, description, part_type, price, stock_quantity, reserved, properties, created_at`

type repository struct {
	pool *pgxpool.Pool

	// getter извлекает активную транзакцию (pgx.Tx) из context.Context
	getter *trmpgx.CtxGetter
}

// NewRepository создаёт репозиторий деталей.
func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}
