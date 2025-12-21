package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-image-web/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BoardRepo struct {
	db *sqlx.DB
}

func NewBoardRepo(db *sqlx.DB) *BoardRepo {
	return &BoardRepo{db: db}
}

func (r *BoardRepo) Create(ctx context.Context, p models.BoardParams) (*models.Board, error) {
	const op = "repo.board.Create"

	b := models.Board{
		Slug: p.Slug,
		Name: p.Name,
		UUID: uuid.NewString(),
	}

	res, err := r.db.NamedExecContext(ctx, `
		INSERT INTO boards (slug, name, uuid)
		VALUES (:slug, :name, :uuid)
	`, &b)
	if err != nil {
		return nil, fmt.Errorf("%s: insert: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: last_insert_id: %w", op, err)
	}

	var out models.Board
	if err := r.db.GetContext(ctx, &out, `
		SELECT id, created_at, slug, name, uuid
		FROM boards
		WHERE id = ?
	`, id); err != nil {
		return nil, fmt.Errorf("%s: select: %w", op, err)
	}

	return &out, nil
}

// dangerous get all
func (r *BoardRepo) GetAll(ctx context.Context) ([]*models.Board, error) {
	const op = "repo.board.GetAll"

	var out []*models.Board
	if err := r.db.SelectContext(ctx, &out, `
    	SELECT id, created_at, slug, name, uuid
    	FROM boards;
	`); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return out, nil
}

func (r *BoardRepo) GetByID(ctx context.Context, id int64) (*models.Board, error) {
	const op = "repo.board.GetByID"

	var out models.Board
	if err := r.db.GetContext(ctx, &out, `
		SELECT id, created_at, slug, name, uuid
		FROM boards
		WHERE id = ?
	`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &out, nil
}

func (r *BoardRepo) GetByUUID(ctx context.Context, uuid string) (*models.Board, error) {
	const op = "repo.board.GetByUUID"

	var out models.Board
	if err := r.db.GetContext(ctx, &out, `
		SELECT id, created_at, slug, name, uuid
		FROM boards
		WHERE uuid = ?
	`, uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &out, nil
}

func (r *BoardRepo) DeleteByID(ctx context.Context, id int64) error {
	const op = "repo.board.DeleteByID"

	res, err := r.db.ExecContext(ctx, `
		DELETE FROM boards
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

func (r *BoardRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	const op = "repo.board.DeleteByUUID"

	res, err := r.db.ExecContext(ctx, `
		DELETE FROM boards
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

// func (r *BoardRepo) UpdateBoard() {}
