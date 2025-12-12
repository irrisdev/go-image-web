package store

import (
	"go-image-web/models"
	"sync"
)

var (
	PostIndex   = map[string]*models.PostUploadModel{}
	PostIndexMu sync.RWMutex
)

func init() {

}
