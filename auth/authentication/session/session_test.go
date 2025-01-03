package session

import (
	"fmt"
	"testing"
	"time"
)

func TestSessionManager(t *testing.T) {

	backend := NewCacheSessionBackend(5, 100, 3*time.Minute)

	sessionManager := NewSessionManager(backend, WithCreateHook(func(sessionID string) {
		fmt.Println("Creating session:", sessionID)
	}))

	// Create session
	sessionID, err := sessionManager.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Created session:", sessionID)

	// Get session
	session, err := sessionManager.GetSession(sessionID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Session data:", session.Data)
	}

	// Update session
	err = sessionManager.UpdateSession(sessionID, "username", "john_doe")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		session, _ = sessionManager.GetSession(sessionID)
		fmt.Println("Updated session data:", session.Data)
	}

	// 更新会话过期时间
	_, err = sessionManager.GetAndRenewSession(sessionID)
	if err != nil {
		fmt.Println("Error touching session:", err)
	} else {
		fmt.Println("Session expiration updated")
	}

	// Destroy session
	err = sessionManager.DestroySession(sessionID)
	if err != nil {
		fmt.Println("Error destroying session:", err)
	} else {
		fmt.Println("Session destroyed")
	}
}
