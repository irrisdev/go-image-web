package store

import (
	"sync"
)

type ThreadUploadState int

const (
	Created ThreadUploadState = iota
	Processing
	Succeeded
	Failed
)

type ThreadUploadStateStore struct {
	mu     sync.RWMutex
	states map[string]ThreadUploadState
}

func (s *ThreadUploadStateStore) Get(token string) (ThreadUploadState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.states[token]
	return state, ok
}

func (s *ThreadUploadStateStore) Set(token string, state ThreadUploadState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[token] = state
}

func NewUploadStateStore() *ThreadUploadStateStore {
	store := &ThreadUploadStateStore{
		states: make(map[string]ThreadUploadState),
	}

	return store

}

// func init() {
// 	// Session cleanup goroutine (optional - clean up old sessions)
// 	go func() {
// 		for {
// 			time.Sleep(30 * time.Minute)
// 			cleanupStale()
// 		}
// 	}()
// }

// // cleanupStale removes tokens older than 1 hour
// func cleanupStale() {
// 	sessionsMu.Lock()
// 	defer sessionsMu.Unlock()

// 	cutoff := time.Now().Add(-2 * time.Hour)
// 	for id, sess := range sessions {
// 		if sess.LastAccess.Before(cutoff) {
// 			delete(sessions, id)
// 		}
// 	}
// }
