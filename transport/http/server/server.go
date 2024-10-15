package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/yates-z/easel/logger"
)

type ServerOption func(*Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

type Server struct {
	*http.Server
	*Router
	listener net.Listener
	network  string
	address  string
	tlsConf  *tls.Config
}

func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		Router:  NewRouter(),
		network: "tcp",
		address: ":80",
	}
	for _, o := range opts {
		o(server)
	}
	server.Server = &http.Server{
		TLSConfig: server.tlsConf,
		Handler:   server.mux,
	}
	return server
}

func (s *Server) Run(ctx context.Context) error {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.listener = listener
		logger.Infof("[http] server listening on: %s", s.listener.Addr().String())
	}
	return s.Serve(s.listener)
}

func (s *Server) MustRun(ctx context.Context) {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			logger.Fatal(err)
		}
		s.listener = listener
		logger.Infof("[http] server listening on: %s", s.listener.Addr().String())
	}
	logger.Fatal(s.Serve(s.listener))
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	logger.Info("[HTTP] server stopping")
	return s.Shutdown(ctx)
}
