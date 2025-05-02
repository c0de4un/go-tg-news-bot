package models

import "time"

type ChatModel struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	ChatID    int64     `db:"chat_id"`
	State     int64     `db:"state"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
