package services

import (
	"go-image-web/models"
	"go-image-web/store"
	"sort"
)

func IndexService() models.IndexPageModel {
	m := models.IndexPageModel{}

	metadata := store.GetImageMetadata()
	for uuid, meta := range metadata {
		// only include images that have variants
		if meta.GetVarientLen() > 0 {
			m.Images = append(m.Images, models.ImageModel{
				ID:        uuid,
				Path:      meta.Original,
				Timestamp: meta.ModifiedTime,
				Size:      meta.FileSize,
			})
		}
	}

	sort.Slice(m.Images, func(i, j int) bool {
		return m.Images[i].Timestamp.Before(m.Images[j].Timestamp)
	})

	return m
}
