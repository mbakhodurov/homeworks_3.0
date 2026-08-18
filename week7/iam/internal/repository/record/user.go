package record

import "time"

// User — persistence-shape пользователя (строка таблицы users).
type User struct {
	UUID         string     `db:"uuid"`
	Login        string     `db:"login"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
