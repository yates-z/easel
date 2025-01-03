package session

import (
	"time"

	"github.com/yates-z/easel/core/cache"
)

// SessionBackend defines a generic interface for session storage.
type SessionBackend interface {
	// Save a session.
	Save(session *Session, ttl time.Duration) error
	// Load a session.
	Load(sessionID string) (*Session, error)
	// Delete a session.
	Delete(sessionID string) error
	// Check if a session exists.
	Exists(sessionID string) bool
}

var _ SessionBackend = (*CacheSessionBackend)(nil)

// CacheSessionBackend implements SessionBackend using a cache backend.
type CacheSessionBackend struct {
	cache cache.Cache[string, *Session]
}

// NewCacheSessionBackend creates a new CacheSessionBackend instance.
func NewCacheSessionBackend(numShards, capacity int, cleanupInterval time.Duration) *CacheSessionBackend {
	return &CacheSessionBackend{
		cache: cache.NewMemCache[string, *Session](numShards, capacity, cleanupInterval),
	}
}

func (s *CacheSessionBackend) Save(session *Session, ttl time.Duration) error {
	if session == nil {
		return ErrSessionIsNil
	}

	return s.cache.Set(session.ID, session, ttl)
}

func (s *CacheSessionBackend) Load(sessionID string) (*Session, error) {
	session, ok := s.cache.Get(sessionID)
	if !ok || session == nil {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (s *CacheSessionBackend) Delete(sessionID string) error {
	s.cache.Delete(sessionID)
	return nil
}

func (s *CacheSessionBackend) Exists(sessionID string) bool {
	return s.cache.HasKey(sessionID)
}
