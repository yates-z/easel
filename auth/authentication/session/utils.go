package session

import (
	"crypto/rand"
	"encoding/hex"
)

func defaultSessionIDGenerator() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("failed to generate session ID")
	}
	return hex.EncodeToString(bytes)
}
