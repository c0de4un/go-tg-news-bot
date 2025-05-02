package models

import "time"

const (
	POST_STATUS_DRAFT      = 0
	POST_STATUS_MODERATION = 1
	POST_STATUS_PUBLISHED  = 2
)

type PostModel struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	Status      int64     `db:"status"`
	UUID        string    `db:"uuid"`
	CreatedBy   int64     `db:"created_by"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   time.Time `db:"deleted_at"`
}
