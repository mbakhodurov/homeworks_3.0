package part

import (
	"context"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/repository/converter"
)

// Create сохраняет новую деталь.
func (r *repository) Create(_ context.Context, part model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[part.UUID] = converter.PartToRecord(part)

	return nil
}
