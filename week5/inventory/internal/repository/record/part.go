package record

import "time"

// Part — плоская структура для маппинга строки таблицы parts.
type Part struct {
	UUID          string    `db:"uuid"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	PartType      string    `db:"part_type"`
	Price         int64     `db:"price"`
	StockQuantity int       `db:"stock_quantity"`
	Reserved      int       `db:"reserved"`
	Properties    []byte    `db:"properties"`
	CreatedAt     time.Time `db:"created_at"`
}
