package session

import "errors"

var (
	ErrSessionIsNil     = errors.New("session cannot be nil")
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionExpired   = errors.New("session expired")
	ErrInvalidSessionID = errors.New("invalid session ID")
)
