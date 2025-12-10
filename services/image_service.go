package services

import (
	"go-image-web/models"
	"go-image-web/store"
	"mime/multipart"
)

type UploadImageResponse struct {
	ID            int
	ErrorMessages []string
}

func GetImages() []models.ImageModel {
	var res = store.GetInstanceImages()

	return res
}

func AddImage(file multipart.File, filename string) (*UploadImageResponse, error) {
	return nil, nil
}
