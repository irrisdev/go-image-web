package models

import (
	"database/sql"
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
