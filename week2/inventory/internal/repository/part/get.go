package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/mbakhodurov/homeworks2/week2/inventory/internal/errors"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/model"
	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/repository/converter"
)

// Get возвращает деталь по UUID.
func (r *repository) Get(_ context.Context, partUUID uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.data[partUUID]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return converter.PartToModel(rec), nil
}
