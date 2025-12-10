package services

import (
	"go-image-web/models"
)

func IndexService() models.IndexPageModel {
	m := models.IndexPageModel{
		Images: GetImages(),
	}

	return m
}
