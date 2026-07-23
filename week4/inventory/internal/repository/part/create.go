package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week4/inventory/internal/model"
)

// Create сохраняет новую деталь. properties и reserved получают значения по умолчанию из схемы БД.
func (r *repository) Create(ctx context.Context, p model.Part) error {
	const query = `INSERT INTO parts (uuid, name, description, part_type, price, stock_quantity, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		p.UUID(), p.Name(), p.Description(), string(p.PartType()), p.Price(), p.StockQuantity(), p.CreatedAt())
	if err != nil {
		return fmt.Errorf("создать деталь: %w", err)
	}

	return nil
}
