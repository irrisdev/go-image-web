package services

import (
	"context"
	"errors"
	"go-image-web/internal/models"
	"go-image-web/internal/repo"
	"log"
)

type BoardService struct {
	repo *repo.BoardRepo
}

func NewBoardService(repo *repo.BoardRepo) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

// create method
var ErrBoardDbError = errors.New("database error when creating new board")

func (s *BoardService) Create(ctx context.Context, p models.BoardParams) (*models.Board, error) {

	// quick validations
	if err := ValidateSlug(p.Slug); err != nil {
		return nil, err
	}

	if err := ValidateName(p.Name); err != nil {
		return nil, err
	}

	board, err := s.repo.Create(ctx, p)
	if err != nil || board == nil {
		log.Println(err)
		return nil, ErrBoardDbError
	}

	return board, nil
}

var (
	ErrSlugLength = errors.New("slug must be less than 6 characters")
	ErrSlugChars  = errors.New("slug must contain letters only (a-z), no digits/spaces/symbols")

	ErrNameLength = errors.New("name must be less than 20 characters")
	ErrNameChars  = errors.New("name must contain letters and spaces only (A-Z/a-z)")
)

func ValidateSlug(s string) error {

	if len(s) > 5 {
		return ErrSlugLength
	}

	for _, v := range s {
		if v < 'a' || v > 'z' {
			return ErrSlugChars
		}
	}

	return nil
}

func ValidateName(s string) error {
	if len(s) > 20 {
		return ErrNameLength
	}

	for _, c := range s {
		if c == ' ' {
			continue
		}

		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
			return ErrNameChars
		}
	}

	return nil
}

// get methods

func (s *BoardService) GetAll(ctx context.Context) ([]*models.Board, error) {
	boards, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (s *BoardService) GetByID(ctx context.Context, id int64) (*models.Board, error) {
	return nil, nil
}

func (s *BoardService) GetByUUID(ctx context.Context, uuid string) (*models.Board, error) {
	return nil, nil
}

func (s *BoardService) GetBySlug(ctx context.Context, slug string) (*models.Board, error) {
	return nil, nil
}

// get all threads for board
func (s *BoardService) List(ctx context.Context, p models.BoardThreadsParams) ([]models.Board, error) {
	return nil, nil
}
