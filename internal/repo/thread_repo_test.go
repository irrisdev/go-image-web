package repo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"go-image-web/internal/models"
	"go-image-web/internal/repo"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupThreadTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite3", "file::memory:?cache=shared&_foreign_keys=on")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Minimal boards table + threads table (FK references boards)
	_, err = db.ExecContext(ctx, `
		CREATE TABLE boards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			slug TEXT NOT NULL UNIQUE COLLATE NOCASE,
			name TEXT NOT NULL,
			uuid TEXT NOT NULL UNIQUE
		);

		CREATE TABLE threads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			uuid TEXT UNIQUE NOT NULL,
			author TEXT NOT NULL,
			subject TEXT NOT NULL,
			message TEXT NOT NULL,

			board_id INTEGER,
			FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE SET NULL
		);

		CREATE INDEX threads_board_id_idx ON threads(board_id);
	`)
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}

	return db
}

func createBoard(t *testing.T, db *sqlx.DB) *models.Board {
	t.Helper()

	br := repo.NewBoardRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	b, err := br.Create(ctx, models.BoardParams{Slug: "g", Name: "Go"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	return b
}

func TestThreadRepo_Create_GetByID_GetByUUID_WithBoard(t *testing.T) {
	db := setupThreadTestDB(t)
	tr := repo.NewThreadRepo(db)

	b := createBoard(t, db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := models.ThreadParams{
		UUID:    "t-uuid-1",
		Author:  "alice",
		Subject: "hello",
		Message: "first post",
		BoardID: b.ID,
	}

	created, err := tr.Create(ctx, in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
	if created.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}
	if created.UUID != in.UUID {
		t.Fatalf("expected uuid %q, got %q", in.UUID, created.UUID)
	}
	if created.Author != in.Author {
		t.Fatalf("expected author %q, got %q", in.Author, created.Author)
	}
	if created.Subject != in.Subject {
		t.Fatalf("expected subject %q, got %q", in.Subject, created.Subject)
	}

	if created.Message != in.Message {
		t.Fatalf("expected message %q, got %q (repo likely sets Message from Subject)", in.Message, created.Message)
	}

	if !created.BoardID.Valid || created.BoardID.Int64 != b.ID {
		t.Fatalf("expected board_id valid=true and %d, got %+v", b.ID, created.BoardID)
	}

	gotByID, err := tr.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if gotByID.ID != created.ID || gotByID.UUID != created.UUID {
		t.Fatalf("GetByID mismatch: got %+v, want %+v", gotByID, created)
	}

	gotByUUID, err := tr.GetByUUID(ctx, created.UUID)
	if err != nil {
		t.Fatalf("GetByUUID: %v", err)
	}
	if gotByUUID.ID != created.ID || gotByUUID.UUID != created.UUID {
		t.Fatalf("GetByUUID mismatch: got %+v, want %+v", gotByUUID, created)
	}
}

func TestThreadRepo_Create_NoBoardID_ShouldBeNull(t *testing.T) {
	db := setupThreadTestDB(t)
	tr := repo.NewThreadRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := models.ThreadParams{
		UUID:    "t-uuid-2",
		Author:  "bob",
		Subject: "no board",
		Message: "hi",
		BoardID: 0, // should be stored as NULL
	}

	created, err := tr.Create(ctx, in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.BoardID.Valid {
		t.Fatalf("expected board_id NULL (Valid=false), got %+v", created.BoardID)
	}
}

func TestThreadRepo_Get_NotFound(t *testing.T) {
	db := setupThreadTestDB(t)
	tr := repo.NewThreadRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tr.GetByID(ctx, 9999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	_, err = tr.GetByUUID(ctx, "does-not-exist")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestThreadRepo_DeleteByID(t *testing.T) {
	db := setupThreadTestDB(t)
	tr := repo.NewThreadRepo(db)

	b := createBoard(t, db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := tr.Create(ctx, models.ThreadParams{
		UUID:    "t-del-id",
		Author:  "a",
		Subject: "s",
		Message: "m",
		BoardID: b.ID,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := tr.DeleteByID(ctx, created.ID); err != nil {
		t.Fatalf("DeleteByID: %v", err)
	}

	_, err = tr.GetByID(ctx, created.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got %v", err)
	}

	if err := tr.DeleteByID(ctx, created.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows on second delete, got %v", err)
	}
}

func TestThreadRepo_DeleteByUUID(t *testing.T) {
	db := setupThreadTestDB(t)
	tr := repo.NewThreadRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := tr.Create(ctx, models.ThreadParams{
		UUID:    "t-del-uuid",
		Author:  "a",
		Subject: "s",
		Message: "m",
		BoardID: 0,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := tr.DeleteByUUID(ctx, created.UUID); err != nil {
		t.Fatalf("DeleteByUUID: %v", err)
	}

	_, err = tr.GetByUUID(ctx, created.UUID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got %v", err)
	}

	if err := tr.DeleteByUUID(ctx, created.UUID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows on second delete, got %v", err)
	}
}
