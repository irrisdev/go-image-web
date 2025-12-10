package store

import (
	"fmt"
	"go-image-web/models"
	"log"
	"os"
	"sync"
)

var ImageAssetsFolder string = "public/assets/images"

var (
	instanceImages []models.ImageModel
	imageMu        sync.RWMutex
)

func init() {
	// load images from filesystem once
	loadImagesFS()
}

func GetInstanceImages() []models.ImageModel {
	imageMu.RLock()
	defer imageMu.RUnlock()
	return instanceImages
}

func AddInstanceImage(image models.ImageModel) models.ImageModel {
	imageMu.Lock()
	defer imageMu.Unlock()
	instanceImages = append(instanceImages, image)
	return instanceImages[len(instanceImages)-1]
}

func loadImagesFS() {
	log.Printf("loading image data from system: %s", ImageAssetsFolder)

	files, err := os.ReadDir(ImageAssetsFolder)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range files {
		fi, err := v.Info()
		if err != nil {
			log.Fatal(err)
		}
		AddInstanceImage(models.ImageModel{
			ID:       i,
			Name:     fi.Name(),
			FilePath: fmt.Sprintf("%s/%s", ImageAssetsFolder, fi.Name()),
			Size:     fi.Size(),
		})
	}
}
