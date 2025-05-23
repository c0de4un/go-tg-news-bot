package models

import "time"

type ForwardPostModel struct {
	ID         int64     `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	FromChatID int64     `json:"from_chat_id" db:"from_chat_id"`
	Status     int       `json:"status" db:"status"`
	CreatedBy  int64     `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
