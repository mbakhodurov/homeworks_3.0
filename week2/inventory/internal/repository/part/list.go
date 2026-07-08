package part

import (
	"context"
	"sort"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week2/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/repository/converter"
)

// List возвращает детали, отфильтрованные по uuids (приоритет) или по типу.
func (r *repository) List(_ context.Context, uuids []uuid.UUID, partType model.PartType) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []model.Part

	if len(uuids) > 0 {
		for _, partUUID := range uuids {
			rec, ok := r.data[partUUID]
			if !ok {
				return nil, errs.ErrPartNotFound
			}
			parts = append(parts, converter.PartToModel(rec))
		}
		return parts, nil
	}

	for _, rec := range r.data {
		if partType == model.PartTypeUnspecified || rec.PartType == string(partType) {
			parts = append(parts, converter.PartToModel(rec))
		}
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}
