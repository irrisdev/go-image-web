package models

import (
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
	Size      int64
	Timestamp time.Time
}

type ImageVarient struct {
	Width int
	Path  string
}

type ImageMetadata struct {
	UUID         string
	OriginalExt  string
	Original     string
	ModifiedTime time.Time
	FileSize     int64

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
