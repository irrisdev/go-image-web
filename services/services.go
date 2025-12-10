package services

import (
	"fmt"
	"go-image-web/models"
	"go-image-web/store"
	"log"
	"os"
)

func IndexService() models.IndexPageModel {
	m := models.IndexPageModel{}

	files, err := os.ReadDir(store.ImageAssetsFolder)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range files {
		fi, err := v.Info()
		if err != nil {
			log.Fatal(err)
		}
		m.Images = append(m.Images, models.ImageModel{
			ID:       i,
			Name:     fi.Name(),
			FilePath: fmt.Sprintf("%s/%s", store.ImageAssetsFolder, fi.Name()),
			Size:     fi.Size(),
		})
	}

	return m
}
