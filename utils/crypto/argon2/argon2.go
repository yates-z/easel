package argon2

import (
	"bytes"
	"crypto/rand"
	"runtime"

	"golang.org/x/crypto/argon2"
)

const (
	memory  = 64 * 1024
	time    = 3
	keyLen  = 32
	saltLen = 16
)

var threads = uint8(runtime.NumCPU())

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password []byte) ([]byte, error) {
	salt, err := generateSalt(saltLen)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(password, salt, time, memory, threads, keyLen)
	hashWithSalt := append(salt, hash...)
	return hashWithSalt, nil
}

func VerifyPassword(password, hash []byte) bool {
	salt := hash[:saltLen]
	expectedHash := hash[saltLen:]
	newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	return bytes.Equal(newHash, expectedHash)
}
