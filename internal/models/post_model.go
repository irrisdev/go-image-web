package models

import (
	"strings"
	"time"
)

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

// NeedsExpand returns true if the message is long enough to need a "View more" toggle.
// More tolerant threshold - only expand if content would exceed image height (~8 lines or 500 chars).
func (p *PostModel) NeedsExpand() bool {
	if len(p.Message) > 500 {
		return true
	}
	newlineCount := strings.Count(p.Message, "\n")
	return newlineCount > 8
}

// TruncatedMessage returns a truncated version of the message ending on a full word.
// Used for the preview when NeedsExpand() is true.
func (p *PostModel) TruncatedMessage() string {
	maxLen := 500
	maxLines := 8

	msg := p.Message

	// First, limit by newlines
	lines := strings.SplitN(msg, "\n", maxLines+1)
	if len(lines) > maxLines {
		msg = strings.Join(lines[:maxLines], "\n")
	}

	// Then, limit by character count and end on a full word
	if len(msg) > maxLen {
		msg = msg[:maxLen]
		// Find the last space to end on a full word
		lastSpace := strings.LastIndexAny(msg, " \t\n")
		if lastSpace > maxLen/2 { // Only truncate to word if we don't lose too much
			msg = msg[:lastSpace]
		}
	}

	return strings.TrimRight(msg, " \t\n")
}
