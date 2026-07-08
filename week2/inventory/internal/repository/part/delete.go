package part

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/mbakhodurov/homeworks2/week2/inventory/internal/errors"
)

// Delete удаляет деталь по UUID.
func (r *repository) Delete(_ context.Context, partUUID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, ok := r.data[partUUID]
	if !ok {
		return errs.ErrPartNotFound
	}

	data.DeletedAt = lo.ToPtr(time.Now())

	return nil
}
