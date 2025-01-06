package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// SessionIDGenerator defines a function for generating session IDs.
type SessionIDGenerator func() string

// SessionManager manages sessions.
type SessionManager struct {
	// Session storage backend
	backend SessionBackend
	// Default session TTL
	sessionTTL time.Duration
	// Session ID generator function
	idGenerator SessionIDGenerator
	// Secret key for HMAC signature
	secretKey []byte
	// Hook for session creation
	onCreate func(sessionID string)
	// Hook for session update
	onUpdate func(sessionID string, key string)
	// Hook for session deletion
	onDestroy func(sessionID string)
}

// SessionManagerOption defines a configuration option for SessionManager.
type SessionManagerOption func(*SessionManager)

// WithIDGenerator sets a custom session ID generator.
func WithIDGenerator(generator SessionIDGenerator) SessionManagerOption {
	return func(sm *SessionManager) {
		sm.idGenerator = generator
	}
}

// WithSecretKey sets a custom secret key for HMAC signature.
func WithSecretKey(secretKey string) SessionManagerOption {
	return func(sm *SessionManager) {
		sm.secretKey = []byte(secretKey)
	}
}

// WithTTL sets a custom TTL for sessions.
func WithTTL(ttl time.Duration) SessionManagerOption {
	return func(sm *SessionManager) {
		sm.sessionTTL = ttl
	}
}

// WithCreateHook sets the on-create hook.
func WithCreateHook(hook func(sessionID string)) SessionManagerOption {
	return func(sm *SessionManager) {
		sm.onCreate = hook
	}
}

// WithUpdateHook sets the on-update hook.
func WithUpdateHook(hook func(sessionID string, key string)) SessionManagerOption {
	return func(sm *SessionManager) {
		sm.onUpdate = hook
	}
}

// SetOnDestroy sets the on-destroy hook.
func WithDestroyHook(hook func(sessionID string)) SessionManagerOption {
	return func(sm *SessionManager) {
		sm.onDestroy = hook
	}
}

// NewSessionManager creates a new SessionManager instance.
func NewSessionManager(backend SessionBackend, opts ...SessionManagerOption) *SessionManager {
	sm := &SessionManager{
		backend:     backend,
		sessionTTL:  30 * time.Minute,
		idGenerator: defaultSessionIDGenerator,
		secretKey:   []byte("default-secret-key"),
	}
	for _, opt := range opts {
		opt(sm)
	}
	return sm
}

// CreateSession generates a new session with a signed ID.
func (sm *SessionManager) CreateSession() (string, error) {
	sessionID := sm.generateSignedSessionID()
	expireAt := time.Now().Add(sm.sessionTTL)
	session := &Session{
		ID:        sessionID,
		ExpiresAt: expireAt,
		Data:      make(map[string]interface{}),
	}
	err := sm.backend.Save(session, sm.sessionTTL)
	if err != nil {
		return "", err
	}
	if sm.onCreate != nil {
		sm.onCreate(sessionID)
	}
	return sessionID, nil
}

// generateSignedSessionID creates a session ID with an HMAC signature.
func (sm *SessionManager) generateSignedSessionID() string {
	sessionID := sm.idGenerator()
	h := hmac.New(sha256.New, sm.secretKey)
	h.Write([]byte(sessionID))
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s.%s", sessionID, signature)
}

// VerifySessionID verifies the HMAC signature of a session ID.
func (sm *SessionManager) VerifySessionID(signedID string) error {
	parts := strings.SplitN(signedID, ".", 2)
	if len(parts) != 2 {
		return ErrInvalidSessionID
	}
	sessionID, signature := parts[0], parts[1]
	h := hmac.New(sha256.New, sm.secretKey)
	h.Write([]byte(sessionID))
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return ErrInvalidSessionID
	}
	return nil
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, error) {
	// Verify the session ID
	if err := sm.VerifySessionID(sessionID); err != nil {
		return nil, err
	}

	session, err := sm.backend.Load(sessionID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (sm *SessionManager) UpdateSession(sessionID string, key string, value interface{}) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}
	session.Data[key] = value
	session.ExpiresAt = time.Now().Add(sm.sessionTTL)
	if sm.onUpdate != nil {
		sm.onUpdate(sessionID, key)
	}
	return sm.backend.Save(session, sm.sessionTTL)
}

func (sm *SessionManager) DestroySession(sessionID string) error {
	// Verify the session ID
	if err := sm.VerifySessionID(sessionID); err != nil {
		return err
	}
	if sm.onDestroy != nil {
		sm.onDestroy(sessionID)
	}
	return sm.backend.Delete(sessionID)
}

func (sm *SessionManager) GetAndRenewSession(sessionID string) (*Session, error) {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	session.ExpiresAt = time.Now().Add(sm.sessionTTL)
	sm.backend.Save(session, sm.sessionTTL)
	return session, nil
}

/*
	RotateSession regenerates the session ID while keeping session data intact.

Which scenarios are suitable for using session rotation mechanism?
After user authentication
When the user successfully logs in, a new session ID is generated to prevent attackers from
hijacking the old session ID before authentication.

Before and after sensitive operations
Rotating session IDs enhances security after users perform certain sensitive operations, such
as changing passwords or making payments.

Regular rotation
Regularly rotate session IDs to reduce the risk of session IDs being brute force cracked or stolen.
*/
func (sm *SessionManager) RotateSession(oldSessionID string) (string, error) {
	// Load the existing session
	sess, err := sm.GetSession(oldSessionID)
	if err != nil {
		return "", err
	}

	// Create a new session ID
	newSessionID := sm.generateSignedSessionID()

	// Save the session data with the new session ID
	if err := sm.backend.Save(sess, sm.sessionTTL); err != nil {
		return "", err
	}

	// Delete the old session
	sm.backend.Delete(oldSessionID)

	return newSessionID, nil
}
