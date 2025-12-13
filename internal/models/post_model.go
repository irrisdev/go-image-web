package models

import "time"

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
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Subject   string    `db:"subject"`
	Message   string    `db:"message"`
	ImageUUID string    `db:"image_uuid"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *PostModel) FormattedTime() string {
	day := p.CreatedAt.Day()
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return p.CreatedAt.Format("Mon, 2") + suffix + p.CreatedAt.Format(" January 2006 15:04")
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
