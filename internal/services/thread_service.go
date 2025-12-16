package services

import "go-image-web/internal/repo"

type ThreadUploadState int

const (
	Created ThreadUploadState = iota
	Processing
	Succeeded
	Failed
)

var threadUploadStates = map[string]ThreadUploadState{}

type ThreadService struct {
	repo repo.ThreadRepo
}

func (s *ThreadService) Create() {}

func (s *ThreadService) Get() {}

func (s *ThreadService) Delete() {}
