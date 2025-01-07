package jwt

import (
	"crypto"
	"crypto/hmac"
	"sync"
)

var signingMethods = map[string]func() SigningMethod{
	"HS256": NewMethodHS256,
	"HS384": NewMethodHS384,
	"HS512": NewMethodHS512,
}
var signingMethodLock = new(sync.RWMutex)

// SigningMethod is a method for signing a token.
type SigningMethod interface {
	Verify(signingString string, sig []byte, key []byte) error
	Sign(signingString string, key []byte) ([]byte, error)
	Alg() string
}

// RegisterSigningMethod registers a signing method.
func RegisterSigningMethod(alg string, f func() SigningMethod) {
	signingMethodLock.Lock()
	defer signingMethodLock.Unlock()

	signingMethods[alg] = f
}

// GetSigningMethod retrieves a signing method from an "alg" string.
func GetSigningMethod(alg string) (method SigningMethod) {
	signingMethodLock.RLock()
	defer signingMethodLock.RUnlock()

	if methodFunc, ok := signingMethods[alg]; ok {
		method = methodFunc()
	}
	return
}

var _ SigningMethod = (*MethodHMAC)(nil)

type MethodHMAC struct {
	Name string
	Hash crypto.Hash
}

// NewMethodHS256 creates a new HMAC signing method using SHA-256.
func NewMethodHS256() SigningMethod {
	return &MethodHMAC{
		Name: "HS256",
		Hash: crypto.SHA256,
	}
}

// NewMethodHS384 creates a new HMAC signing method using SHA-384.
func NewMethodHS384() SigningMethod {
	return &MethodHMAC{
		Name: "HS384",
		Hash: crypto.SHA384,
	}
}

// NewMethodHS512 creates a new HMAC signing method using SHA-512.
func NewMethodHS512() SigningMethod {
	return &MethodHMAC{
		Name: "HS512",
		Hash: crypto.SHA512,
	}
}

// Sign implements SigningMethod.
func (m *MethodHMAC) Sign(signingString string, key []byte) ([]byte, error) {
	h := hmac.New(m.Hash.New, key)
	h.Write([]byte(signingString))
	return h.Sum(nil), nil
}

// Verify implements SigningMethod.
func (m *MethodHMAC) Verify(signingString string, sig []byte, key []byte) error {
	// Can we use the specified hashing method?
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}

	h := hmac.New(m.Hash.New, key)
	h.Write([]byte(signingString))
	if !hmac.Equal(sig, h.Sum(nil)) {
		return ErrSignatureInvalid
	}
	return nil
}

func (m *MethodHMAC) Alg() string {
	return m.Name
}
