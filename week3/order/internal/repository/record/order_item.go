package record

import "time"

// OrderItem — persistence-shape позиции заказа (строка таблицы order_items).
type OrderItem struct {
	UUID      string    `db:"uuid"`
	OrderUUID string    `db:"order_uuid"`
	PartUUID  string    `db:"part_uuid"`
	PartType  string    `db:"part_type"`
	Price     int64     `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}
