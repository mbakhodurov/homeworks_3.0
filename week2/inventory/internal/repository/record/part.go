package record

import "time"

// Part — persistence-shape детали для хранилища.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      string
	StockQuantity int64
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}
