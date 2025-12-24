package models

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type Thread struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	UUID    string `db:"uuid" json:"uuid"`
	Author  string `db:"author" json:"author"`
	Subject string `db:"subject" json:"subject"`
	Message string `db:"message" json:"message"`

	BoardID sql.NullInt64 `db:"board_id" json:"board_id"`
}

type ThreadParams struct {
	UUID    string
	Author  string
	Subject string
	Message string
	BoardID int64
}

type NewThreadInputs struct {
	File           multipart.File
	Header         *multipart.FileHeader
	Subject        string
	Message        string
	BoardID        int64
	IdempotencyKey string
}

type ThreadView struct {
	Slug   string
	Thread *Thread
	Image  *ImageModel
	Error  string
}

type ThreadItem struct {
	Thread *Thread
	Image  *ImageModel
}

func (t *Thread) FormattedTime() string {
	day := t.CreatedAt.Day()
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.CreatedAt.Format("Mon, 2") + suffix + t.CreatedAt.Format(" January 2006 15:04")
}
