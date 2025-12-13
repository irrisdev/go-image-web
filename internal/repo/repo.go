package repo

import (
	"fmt"
	"go-image-web/internal/models"

	"github.com/jmoiron/sqlx"
)

type PostRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

const allPostsQuery string = `
SELECT posts.id,
       posts.name,
       posts.subject,
       posts.message,
       posts.image_uuid,
	   posts.created_at
FROM posts
ORDER BY posts.created_at DESC;
`

func (r *PostRepo) SelectAllPosts() ([]*models.PostModel, error) {

	var posts []*models.PostModel
	if err := r.db.Select(&posts, allPostsQuery); err != nil {
		return nil, err
	}

	return posts, nil
}

const insertPostQuery string = `
INSERT INTO posts(name, subject, message, image_uuid)
VALUES (:name,
        :subject,
        :message,
        :image_uuid
)
RETURNING id, name, subject, message, image_uuid;
`

func (r *PostRepo) InsertPost(entry *models.PostModel) (*models.PostModel, error) {
	const op string = "repo.post.InsertPost"

	rows, err := r.db.NamedQuery(insertPostQuery, &entry)
	if err != nil {
		return &models.PostModel{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	if rows.Next() {
		var out models.PostModel
		if err := rows.StructScan(&out); err != nil {
			return &models.PostModel{}, fmt.Errorf("%s: scan: %w", op, err)

		}
		return &out, nil

	}
	return &models.PostModel{}, fmt.Errorf("%s: no row returned", op)

}
