package backend

import (
	"os"
)

var _ Backend = (*osBackend)(nil)

type osBackend struct {
}

type OSBackendBuilder struct {
	backend *osBackend
}

func (b *OSBackendBuilder) Build() *osBackend {
	return b.backend
}

func OSBackend() *OSBackendBuilder {
	return &OSBackendBuilder{
		backend: &osBackend{},
	}
}

func (O osBackend) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}

func (O osBackend) Sync() error {
	return nil
}

func (O osBackend) AllowANSI() bool {
	return true
}
