package models

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type IndexPageModel struct {
	Images         []ImageModel
	OriginalImages []string
}

type ImageModel struct {
	ID        string
	Path      string
	Extension string
	Width     int
	Height    int
	Timestamp time.Time
	Size      int64
}

// ShortID returns the first 8 characters of the ID
func (m ImageModel) ShortID() string {
	parts := strings.Split(m.ID, "-")
	if len(parts) < 1 {
		return m.ID
	}
	return parts[0]
}

func (m ImageModel) FormattedTime() string {
	day := m.Timestamp.Day()
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return m.Timestamp.Format("Mon, 2") + suffix + m.Timestamp.Format(" January 2006 15:04")
}

func (m ImageModel) FormattedSize() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case m.Size >= GB:
		return fmt.Sprintf("%.2f GB", float64(m.Size)/float64(GB))
	case m.Size >= MB:
		return fmt.Sprintf("%.2f MB", float64(m.Size)/float64(MB))
	case m.Size >= KB:
		return fmt.Sprintf("%.2f KB", float64(m.Size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", m.Size)
	}
}

type ImageVarient struct {
	Width int
	Path  string
}

type ImageMetadata struct {
	UUID         string
	OriginalExt  string
	OriginalPath string
	ModifiedTime time.Time
	OriginalSize int64

	OriginalWidth  int
	OriginalHeight int

	VarientsMu sync.RWMutex
	Varients   map[int]ImageVarient
}

func (m *ImageMetadata) SetVariant(width int, v ImageVarient) {
	m.VarientsMu.Lock()
	defer m.VarientsMu.Unlock()
	if m.Varients == nil {
		m.Varients = make(map[int]ImageVarient)
	}
	m.Varients[width] = v
}

func (m *ImageMetadata) GetVariant(width int) (ImageVarient, bool) {
	m.VarientsMu.RLock()
	defer m.VarientsMu.RUnlock()
	v, ok := m.Varients[width]
	return v, ok
}

func (m *ImageMetadata) GetVarientLen() int {
	m.VarientsMu.RLock()
	defer m.VarientsMu.RUnlock()
	return len(m.Varients)
}
