package models

type IndexPageModel struct {
	Images []ImageModel
}

type ImageModel struct {
	ID       int // Index of image, assigned based on position in file
	Name     string
	FilePath string
	Size     int64
}
