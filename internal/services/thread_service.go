package services

import (
	"context"
	"go-image-web/internal/models"
	"go-image-web/internal/repo"
	"go-image-web/internal/store"

	"github.com/google/uuid"
)

type ThreadService struct {
	repo         *repo.ThreadRepo
	uploadStates *store.ThreadUploadStateStore
}

func NewThreadService(repo *repo.ThreadRepo) *ThreadService {
	return &ThreadService{
		repo:         repo,
		uploadStates: store.NewUploadStateStore(),
	}
}

func (s *ThreadService) Create() {}

func (s *ThreadService) Get() {}

func (s *ThreadService) Delete() {}

func (s *ThreadService) GetByBoardID(ctx context.Context, id int) ([]*models.Thread, error) {

	threads, err := s.repo.ListByBoardID(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (s *ThreadService) NewUploadToken() string {
	token := uuid.New().String()
	s.uploadStates.Set(token, store.Created)
	return token
}
