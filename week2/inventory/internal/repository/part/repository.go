package part

import (
	"sync"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/inventory/internal/repository/record"
)

// repository — in-memory хранилище деталей.
type repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]record.Part
}

// New создаёт новый in-memory репозиторий деталей с предзагруженными seed-данными.
func NewRepository() *repository {

	return &repository{data: make(map[uuid.UUID]record.Part)}
}
