package session

import (
	"sync/atomic"
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

	Times atomic.Int32

	DeviceType int8
}
