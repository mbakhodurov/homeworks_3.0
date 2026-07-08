package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/mbakhodurov/homeworks2/week2/order/internal/repository/record"
)

// repository — in-memory хранилище заказов.
type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

// New создаёт новый in-memory репозиторий заказов.
func New() *repository {
	return &repository{orders: make(map[uuid.UUID]record.Order)}
}
