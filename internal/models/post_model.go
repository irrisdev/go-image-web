package models

// POST upload binding
type PostUploadModel struct {
	Name    string
	Subject string
	Message string
}

type PostViewModel struct {
	Post  *PostModel
	Image *ImageModel
}

type PostModel struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Subject   string `db:"subject"`
	Message   string `db:"message"`
	ImageUUID string `db:"image_uuid"`
}

func (p *PostModel) NewPost(name string, subject string, msg string, uuid string) *PostModel {

	if name == "" {
		name = "Anon"
	}

	return &PostModel{
		Name:      name,
		Subject:   subject,
		Message:   msg,
		ImageUUID: uuid,
	}
}
