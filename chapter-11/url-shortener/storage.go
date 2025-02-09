package main

import "sync"

// MemoryStore is an in-memory store for URL mappings
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewMemoryStore initializes a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]string),
	}
}

// Save stores a short URL and its corresponding long URL
func (s *MemoryStore) Save(shortCode, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[shortCode] = longURL
}

// Get retrieves the original URL for a given short code
func (s *MemoryStore) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, found := s.store[shortCode]
	return longURL, found
}
