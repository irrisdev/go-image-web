package services

import (
	"fmt"
	"go-image-web/internal/models"
	"go-image-web/internal/repo"
	"log"
)

type PostService struct {
	repo *repo.PostRepo
}

func NewPostService(repo *repo.PostRepo) *PostService {
	return &PostService{
		repo: repo,
	}
}

const (
	DefaultPagePosts int    = 5
	DefaultPostName  string = "Anon"
)

func (p *PostService) GetPosts() (map[int]*models.PostModel, error) {
	posts, err := p.repo.SelectAllPosts()
	if err != nil {
		log.Printf("failed to select all posts %v", err)
		return nil, fmt.Errorf("failed to retireve posts")
	}

	var postMap = make(map[int]*models.PostModel)
	for _, v := range posts {
		postMap[v.ID] = &models.PostModel{
			ID:        v.ID,
			Name:      v.Name,
			Subject:   v.Subject,
			Message:   v.Message,
			ImageUUID: v.ImageUUID,
		}
	}

	return postMap, nil
}

func (p *PostService) SavePost(model *models.PostModel) (int, error) {

	if model == nil {
		return 0, fmt.Errorf("nil reference passed to SavePost")
	}

	createdModel, err := p.repo.InsertPost(model)
	if err != nil {
		return 0, err
	}

	return createdModel.ID, nil
}
