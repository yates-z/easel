package backend

import (
	"os"
)

var _ Backend = (*OsBackend)(nil)

type OsBackend struct {
	WriteSyncer
}

type OSBackendBuilder struct {
	backend *OsBackend
}

func (b *OSBackendBuilder) Build() *OsBackend {
	return b.backend
}

func OSBackend() *OSBackendBuilder {
	return &OSBackendBuilder{
		backend: &OsBackend{
			os.Stderr,
		},
	}
}

func (b *OsBackend) Write(p []byte) (n int, err error) {
	return b.WriteSyncer.Write(p)
}

func (b *OsBackend) Sync() error {
	return b.WriteSyncer.Sync()
}

func (b *OsBackend) Close() error {
	return nil
}

func (b *OsBackend) AllowANSI() bool {
	return true
}
