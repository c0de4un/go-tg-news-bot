package models

import "time"

type MediaModel struct {
	ID         int64     `db:"id"`
	FileName   string    `db:"file_name"`
	LocalSrc   string    `db:"local_str"`
	Extension  string    `db:"extension"`
	MimeType   string    `db:"mime_type"`
	Length     int64     `db:"length"`
	UploadedBy int64     `db:"uploaded_by"`
	PostID     int64     `db:"post_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
