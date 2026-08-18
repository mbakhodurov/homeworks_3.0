package part

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mbakhodurov/homeworks2/week7/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week7/inventory/internal/repository/converter"
)

// Create сохраняет новую деталь в БД, включая properties в JSONB.
func (r *repository) Create(ctx context.Context, p model.Part) error {
	propsJSON, err := json.Marshal(converter.PartPropertiesToRecord(*p.Properties()))
	if err != nil {
		return fmt.Errorf("сериализовать свойства: %w", err)
	}

	const query = `INSERT INTO parts (uuid, name, description, part_type, price, stock_quantity, properties, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		p.UUID(), p.Name(), p.Description(), string(p.PartType()), p.Price(), p.StockQuantity(), propsJSON, p.CreatedAt())
	if err != nil {
		return fmt.Errorf("создать деталь: %w", err)
	}

	return nil
}
