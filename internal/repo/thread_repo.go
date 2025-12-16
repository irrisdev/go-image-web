package repo

import "github.com/jmoiron/sqlx"

type ThreadRepo struct {
	db *sqlx.DB
}

func NewThreadRepo(db *sqlx.DB) *ThreadRepo {
	return &ThreadRepo{db: db}
}

func (r *ThreadRepo) Create() {}

func (r *ThreadRepo) Get() {}

func (r *ThreadRepo) Delete() {}
