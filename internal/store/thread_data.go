package store

import (
	"sync"
	"time"
)

type UploadState int

const (
	Created UploadState = iota
	Processing
	Succeeded
	Failed
)

type UploadEntry struct {
	CreatedAt  time.Time
	State      UploadState
	ThreadUUID string
}

type UploadStateStore struct {
	mu     sync.RWMutex
	states map[string]UploadEntry
}

func NewUploadStateStore() *UploadStateStore {
	return &UploadStateStore{
		states: make(map[string]UploadEntry),
	}
}

func (s *UploadStateStore) Get(token string) (UploadEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.states[token]
	return entry, ok
}

func (s *UploadStateStore) Set(token string, state UploadState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[token] = UploadEntry{
		CreatedAt: time.Now(),
		State:     state,
	}
}

func (s *UploadStateStore) Update(token string, state UploadState, uuid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := s.states[token]
	entry.State = state
	entry.ThreadUUID = uuid
	s.states[token] = entry
}

func (s *UploadStateStore) Cleanup(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)

	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.states {
		if v.CreatedAt.Before(cutoff) {
			delete(s.states, k)
		}
	}
}
