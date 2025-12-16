package services

import (
	"context"
	"go-image-web/internal/models"
	"go-image-web/internal/repo"
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

func (s *BoardService) Create(ctx context.Context, p models.BoardParams) (*models.Board, error) {
	return nil, nil
}

// get methods

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
