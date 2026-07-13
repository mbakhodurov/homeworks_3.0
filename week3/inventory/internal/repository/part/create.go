package part

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week3/inventory/internal/repository/converter"
)

// Create сохраняет новую деталь.
func (r *repository) Create(ctx context.Context, part model.Part) error {
	const query = `INSERT INTO parts (uuid, name, description, part_type, price, stock_quantity, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	rec := converter.PartToRecord(part)

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, rec.UUID, rec.Name, rec.Description, rec.PartType, rec.Price, rec.StockQuantity, rec.CreatedAt)
	if err != nil {
		return fmt.Errorf("создать деталь: %w", err)
	}

	return nil
}
