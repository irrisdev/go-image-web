package services

import (
	"go-image-web/internal/models"
	"go-image-web/internal/store"
	"sort"
)

func IndexService() models.IndexPageModel {
	m := models.IndexPageModel{}

	metadata := store.GetImageMetadata()
	for uuid, meta := range metadata {
		if meta.GetVarientLen() > 0 {
			m.Images = append(m.Images, models.ImageModel{
				ID:        uuid,
				Path:      meta.OriginalPath,
				Extension: meta.OriginalExt,
				Width:     meta.OriginalWidth,
				Height:    meta.OriginalHeight,
				Timestamp: meta.ModifiedTime,
				Size:      meta.OriginalSize,
			})
		}
	}

	sort.Slice(m.Images, func(i, j int) bool {
		return m.Images[i].Timestamp.After(m.Images[j].Timestamp)
	})

	return m
}
