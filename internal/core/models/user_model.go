package models

import "time"

type UserModel struct {
	ID               int64     `db:"id"`
	TelegramUsername string    `db:"telegram_username"`
	TelegramID       int64     `db:"telegram_id"`
	UUID             string    `db:"uuid"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
