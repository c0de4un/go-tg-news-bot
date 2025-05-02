package models

import "time"

const (
	CHAT_STATE_CLIENT_REGISTRATION = 1
	CHAT_STATE_POST_TITLE          = 2
	CHAT_STATE_POST_BODY           = 3
	CHAT_STATE_POST_MEDIA          = 4
	CHAT_STATE_POST_WELCOME        = 0
)

type ChatModel struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	ChatID    int64     `db:"chat_id"`
	State     int64     `db:"state"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
