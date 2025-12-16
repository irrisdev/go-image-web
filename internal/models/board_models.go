package models

import "time"

type Board struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	Slug string `db:"slug" json:"slug"`
	Name string `db:"name" json:"name"`
	UUID string `db:"uuid" json:"uuid"`
}

type BoardParams struct {
	Slug string
	Name string
}
