package part

import (
	"context"
	"sort"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/repository/converter"
)

// GetAll возвращает все детали без фильтрации, отсортированные по имени.
func (r *repository) GetAll(_ context.Context) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.data))
	for _, rec := range r.data {
		parts = append(parts, converter.PartToModel(rec))
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}
