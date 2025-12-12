package repo

import "github.com/jmoiron/sqlx"

type PostRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}
