package models

import "time"

const (
	CLIENT_STATE_REGISTRATION = 0
	CLIENT_STATE_DRAFT        = 1
	CLIENT_STATE_ACTIVE       = 2
	CLIENT_STATE_DISABLED     = 3
)

type ClientModel struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	Status    int64      `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
