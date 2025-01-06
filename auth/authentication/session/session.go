package session

import (
	"time"
)

// Session represents a single session
type Session struct {
	// Session ID
	ID string
	// Expiration time
	ExpiresAt time.Time
	// Data stored in the session
	Data map[string]interface{}
}
