package services

import "go-image-web/models"

func IndexService() models.IndexPageModel {
	m := models.IndexPageModel{
		Images: []models.ImageModel{
			{
				ID:       0,
				Name:     "JohnDoe",
				FilePath: "/public/assets/1.jpg",
				Size:     100,
			},
		},
	}
	return m
}
