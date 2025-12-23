package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-image-web/internal/models"

	"github.com/jmoiron/sqlx"
)

type ThreadRepo struct {
	db *sqlx.DB
}

func NewThreadRepo(db *sqlx.DB) *ThreadRepo {
	return &ThreadRepo{db: db}
}

func (r *ThreadRepo) Create(ctx context.Context, p models.ThreadParams) (*models.Thread, error) {
	const op = "repo.thread.Create"

	t := models.Thread{
		UUID:    p.UUID,
		Author:  p.Author,
		Subject: p.Subject,
		Message: p.Message,
		BoardID: sql.NullInt64{Int64: p.BoardID, Valid: p.BoardID > 0},
	}

	res, err := r.db.NamedExecContext(ctx, `
		INSERT INTO threads(uuid, author, subject, message, board_id)
		VALUES (:uuid, :author, :subject, :message, :board_id)
	`, &t)
	if err != nil {
		return nil, fmt.Errorf("%s: insert: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: last_insert_id: %w", op, err)
	}

	var out models.Thread
	if err := r.db.GetContext(ctx, &out, `
		SELECT id, created_at, uuid, author, subject, message, board_id
		FROM threads
		WHERE id = ?
	`, id); err != nil {
		return nil, fmt.Errorf("%s: select: %w", op, err)
	}

	return &out, nil
}

func (r *ThreadRepo) GetByID(ctx context.Context, id int64) (*models.Thread, error) {
	const op = "repo.thread.GetByID"

	var out models.Thread
	if err := r.db.GetContext(ctx, &out, `
		SELECT id, created_at, uuid, author, subject, message, board_id
		FROM threads
		WHERE id = ?
	`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &out, nil
}

func (r *ThreadRepo) GetByUUID(ctx context.Context, uuid string) (*models.Thread, error) {
	const op = "repo.thread.GetByUUID"

	var out models.Thread
	if err := r.db.GetContext(ctx, &out, `
		SELECT id, created_at, uuid, author, subject, message, board_id
		FROM threads
		WHERE uuid = ?
	`, uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &out, nil
}

func (r *ThreadRepo) DeleteByID(ctx context.Context, id int64) error {
	const op = "repo.thread.DeleteByID"

	res, err := r.db.ExecContext(ctx, `
		DELETE FROM threads
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("%s: delete: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows_affected: %w", op, err)
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ThreadRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	const op = "repo.thread.DeleteByUUID"

	res, err := r.db.ExecContext(ctx, `
		DELETE FROM threads
		WHERE uuid = ?
	`, uuid)
	if err != nil {
		return fmt.Errorf("%s: delete: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows_affected: %w", op, err)
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ThreadRepo) ListByBoardID(ctx context.Context, id int64) ([]*models.Thread, error) {

	const op = "repo.thread.ListByBoardID"

	var out []*models.Thread
	if err := r.db.SelectContext(ctx, &out, `
    	SELECT id, created_at, uuid, author, subject, message, board_id
    	FROM threads
		WHERE board_id = ?;
	`, id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil

}
