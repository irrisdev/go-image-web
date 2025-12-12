package models

type PostUploadModel struct {
	Name    string
	Subject string
	Message string
}

type PostViewModel struct {
	Name    string
	Subject string
	Message string
	Image   ImageModel
}
