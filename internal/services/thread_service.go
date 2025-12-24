package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-image-web/internal/models"
	"go-image-web/internal/repo"
	"go-image-web/internal/store"
	"time"

	"github.com/google/uuid"
)

type ThreadService struct {
	repo         *repo.ThreadRepo
	uploadStates *store.UploadStateStore
}

func NewThreadService(repo *repo.ThreadRepo) *ThreadService {
	return &ThreadService{
		repo:         repo,
		uploadStates: store.NewUploadStateStore(),
	}
}

var (
	ErrValidationSubject = errors.New("subject too long")
	ErrValidationMessage = errors.New("message too long")
	ErrMissingImage      = errors.New("must provide an image")
	ErrBadIdempotencyKey = errors.New("bad request, try again")
	ErrThreadExists      = errors.New("thread already uploaded")
	ErrImageTooBig       = errors.New("image too big")
	ErrInsertFailure     = errors.New("failed to create new thread")
	ErrFileError         = errors.New("failed to open file")
)

const (
	MaxThreadBytes  int64 = 15 << 20
	MaxSubjectChars int   = 70
	MaxMessageChars int   = 1500
)

func (s *ThreadService) Create(ctx context.Context, p *models.NewThreadInputs) (int, error) {

	if len(p.Subject) > MaxSubjectChars {
		return 0, ErrValidationSubject
	}

	if len(p.Message) > MaxMessageChars {
		return 0, ErrValidationMessage
	}

	// check if size in header is too big
	if p.Header.Size > MaxThreadBytes {
		return 0, ErrImageTooBig
	}

	uuid := p.IdempotencyKey
	tmpPath, err := store.CreateTmpFile(uuid, p.File)
	p.File.Close()
	if err != nil {
		return 0, ErrFileError
	}

	go func() {
		s.uploadStates.Update(p.IdempotencyKey, store.Processing, uuid)

		_, err := SaveImage(tmpPath, p.Header.Filename, uuid)
		if err != nil {
			s.uploadStates.Update(p.IdempotencyKey, store.Failed, uuid)
			return
		}
		s.uploadStates.Update(p.IdempotencyKey, store.Succeeded, uuid)
	}()

	thread, err := s.repo.Create(ctx, models.ThreadParams{
		UUID:    uuid,
		Author:  "Anonymous",
		Subject: p.Subject,
		Message: p.Message,
		BoardID: p.BoardID,
	})

	if err != nil {
		return 0, ErrInsertFailure
	}

	return int(thread.ID), nil
}

func (s *ThreadService) GetByID(ctx context.Context, id int) (*models.Thread, error) {
	thread, err := s.repo.GetByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("thread %d does not exist", id)
		}
		return nil, fmt.Errorf("internal server error occured")
	}
	return thread, nil
}

func (s *ThreadService) GetImageByUUID(uuid string) *models.ImageModel {

	meta := store.GetGuidImageMetadata(uuid)
	if meta != nil {
		return &models.ImageModel{
			ID:        meta.UUID,
			Path:      meta.OriginalPath,
			Extension: meta.OriginalExt,
			Width:     meta.OriginalWidth,
			Height:    meta.OriginalHeight,
			Timestamp: meta.ModifiedTime,
			Size:      meta.OriginalSize,
		}
	}
	return nil
}

func (s *ThreadService) GetListByBoardID(ctx context.Context, id int) ([]*models.ThreadItem, error) {

	threads, err := s.repo.ListByBoardID(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	var threadItems []*models.ThreadItem
	for _, thread := range threads {
		item := &models.ThreadItem{
			Thread: thread,
		}

		image := s.GetImageByUUID(thread.UUID)
		if image != nil {
			item.Image = image
		}

		threadItems = append(threadItems, item)
	}

	return threadItems, nil
}

func (s *ThreadService) NewUploadToken() string {
	token := uuid.New().String()
	s.uploadStates.Set(token, store.Created)
	return token
}

func (s *ThreadService) GetUploadEntry(key string) (store.UploadEntry, bool) {
	return s.uploadStates.Get(key)
}

func (s *ThreadService) StartStateCleanup(ctx context.Context, every, ttl time.Duration) {
	t := time.NewTicker(every)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.uploadStates.Cleanup(ttl)
			}
		}
	}()
}
