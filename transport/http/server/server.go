package server

import (
	"context"
	"crypto/tls"
	"html/template"
	"net"
	"net/http"

	"github.com/yates-z/easel/logger"
	templ "github.com/yates-z/easel/transport/http/server/template"
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
	return func(s *Server) {
		s.tlsConf = c
	}
}

// Middlewares with global middleware.
func Middlewares(middlewares ...Middleware) ServerOption {
	return func(s *Server) {
		s.Use(middlewares...)
	}
}

// ShowInfo with showInfo config.
func ShowInfo(isShow bool) ServerOption {
	return func(s *Server) {
		s.showInfo = isShow
	}
}

// HTMLTemplate with htmlTemplate config.
func HTMLTemplate(t *templ.HTMLTemplate) ServerOption {
	return func(s *Server) {
		s.htmlTempl = t
	}
}

type Server struct {
	*http.Server
	*Router
	listener net.Listener
	network  string
	address  string
	tlsConf  *tls.Config

	showInfo  bool
	htmlTempl *templ.HTMLTemplate
}

func New(opts ...ServerOption) *Server {
	server := &Server{
		network: "tcp",
		address: ":80",

		showInfo:  false,
		htmlTempl: templ.New(),
	}
	server.Router = NewRouter(server)
	for _, o := range opts {
		o(server)
	}
	server.Server = &http.Server{
		TLSConfig: server.tlsConf,
		Handler:   server.mux,
	}
	return server
}

// Template returns server.htmlTempl.
func (s *Server) Template() *template.Template {
	return s.htmlTempl.Templ
}

// LoadHTMLFiles loads a slice of HTML files.
func (s *Server) LoadHTMLGlob(pattern string) {
	s.htmlTempl.LoadHTMLGlob(pattern)
}

// LoadHTMLGlob loads HTML files identified by glob pattern.
func (s *Server) LoadHTMLFiles(files ...string) {
	s.htmlTempl.LoadHTMLFiles(files...)
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
	if s.tlsConf != nil {
		return s.ServeTLS(s.listener, "", "")
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
	if s.tlsConf != nil {
		logger.Fatal(s.ServeTLS(s.listener, "", ""))
	}
	logger.Fatal(s.Serve(s.listener))
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	logger.Info("[HTTP] server stopping")
	return s.Shutdown(ctx)
}
