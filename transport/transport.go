package transport

import (
	"context"
	"net/url"
)

// Server is transport server.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Endpoint() (*url.URL, error)
}
