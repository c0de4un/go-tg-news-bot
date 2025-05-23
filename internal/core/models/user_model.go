package models

import "time"

const (
	USER_ROLE_CLIENT = 0
	USER_ROLE_GUEST  = 1
	USER_ROLE_ADMIN  = 2
)

type UserModel struct {
	ID               int64     `db:"id"`
	TelegramUsername string    `db:"telegram_username"`
	TelegramID       int64     `db:"telegram_id"`
	UUID             string    `db:"uuid"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	Role             int64     `db:"role"`
	Client           *ClientModel
	Chat             *ChatModel
	ChatID           int64
}
