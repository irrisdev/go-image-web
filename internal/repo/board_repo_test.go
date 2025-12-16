package repo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"go-image-web/internal/models"
	"go-image-web/internal/repo"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	// Use a shared in-memory DB + single connection to keep schema/data stable.
	db, err := sqlx.Open("sqlite3", "file::memory:?cache=shared&_foreign_keys=on")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	schema := `
	CREATE TABLE boards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		slug TEXT NOT NULL UNIQUE COLLATE NOCASE,
		name TEXT NOT NULL,
		uuid TEXT NOT NULL UNIQUE
	);
	`
	if _, err := db.ExecContext(ctx, schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	return db
}

func TestBoardRepo_Create_GetByID_GetByUUID(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewBoardRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := models.BoardParams{
		Slug: "g",
		Name: "Go",
	}

	created, err := r.Create(ctx, in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
	if created.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}
	if created.Slug != in.Slug {
		t.Fatalf("expected slug %q, got %q", in.Slug, created.Slug)
	}
	if created.Name != in.Name {
		t.Fatalf("expected name %q, got %q", in.Name, created.Name)
	}
	if created.UUID == "" {
		t.Fatalf("expected UUID to be set")
	}
	if _, err := uuid.Parse(created.UUID); err != nil {
		t.Fatalf("expected valid UUID, got %q: %v", created.UUID, err)
	}

	gotByID, err := r.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if gotByID.ID != created.ID || gotByID.UUID != created.UUID {
		t.Fatalf("GetByID mismatch: got %+v, want %+v", gotByID, created)
	}

	gotByUUID, err := r.GetByUUID(ctx, created.UUID)
	if err != nil {
		t.Fatalf("GetByUUID: %v", err)
	}
	if gotByUUID.ID != created.ID || gotByUUID.UUID != created.UUID {
		t.Fatalf("GetByUUID mismatch: got %+v, want %+v", gotByUUID, created)
	}
}

func TestBoardRepo_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewBoardRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.GetByID(ctx, 9999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	_, err = r.GetByUUID(ctx, uuid.NewString())
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestBoardRepo_DeleteByID(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewBoardRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := r.Create(ctx, models.BoardParams{Slug: "x", Name: "X"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := r.DeleteByID(ctx, created.ID); err != nil {
		t.Fatalf("DeleteByID: %v", err)
	}

	// Ensure it is gone.
	_, err = r.GetByID(ctx, created.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got %v", err)
	}

	// Deleting again should be not found.
	if err := r.DeleteByID(ctx, created.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows on second delete, got %v", err)
	}
}

func TestBoardRepo_DeleteByUUID(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewBoardRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := r.Create(ctx, models.BoardParams{Slug: "y", Name: "Y"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := r.DeleteByUUID(ctx, created.UUID); err != nil {
		t.Fatalf("DeleteByUUID: %v", err)
	}

	_, err = r.GetByUUID(ctx, created.UUID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got %v", err)
	}

	if err := r.DeleteByUUID(ctx, created.UUID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows on second delete, got %v", err)
	}
}
