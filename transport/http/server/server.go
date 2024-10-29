package server

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/yates-z/easel/core/pool"
	"github.com/yates-z/easel/logger"
	templ "github.com/yates-z/easel/transport/http/server/template"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
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
func ShowInfo() ServerOption {
	return func(s *Server) {
		s.showInfo = true
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

	ctxPool      *pool.Pool[*Context]
	middlewares  []Middleware
	showInfo     bool
	htmlTempl    *templ.HTMLTemplate
	errorHandler func(ctx *Context, err error)
}

func New(opts ...ServerOption) *Server {
	server := &Server{
		network: "tcp",
		address: ":80",

		showInfo:  false,
		htmlTempl: templ.New(),
	}
	server.Router = NewRouter(server)
	server.ctxPool = pool.New(func() *Context {
		return newContext(nil, server)
	})
	server.errorHandler = func(ctx *Context, err error) {
		http.Error(ctx.Response, err.Error(), http.StatusBadRequest)
	}
	for _, o := range opts {
		o(server)
	}

	handler := chain(server.middlewares...)(func(ctx *Context) error {
		server.mux.ServeHTTP(ctx.Response, ctx.Request)
		return nil
	})

	server.Server = &http.Server{
		TLSConfig: server.tlsConf,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			ctx := server.ctxPool.Get()
			ctx.WithBaseContext(req.Context())
			req = req.Clone(ctx)
			ctx.init(req, resp)
			if err := handler(ctx); err != nil {
				server.errorHandler(ctx, err)
			}
			ctx.reset()
			server.ctxPool.Put(ctx)
		}),
	}
	return server
}

func (s *Server) Use(middlewares ...Middleware) *Server {
	s.middlewares = append(s.middlewares, middlewares...)
	return s
}

// Template returns server.htmlTempl.
func (s *Server) Template() *template.Template {
	return s.htmlTempl.Tpl
}

// LoadHTMLGlob loads a slice of HTML files.
func (s *Server) LoadHTMLGlob(pattern string) {
	s.htmlTempl.LoadHTMLGlob(pattern)
}

// LoadHTMLFiles loads HTML files identified by glob pattern.
func (s *Server) LoadHTMLFiles(files ...string) {
	s.htmlTempl.LoadHTMLFiles(files...)
}

// SetErrorHandler sets custom http error handler.
func (s *Server) SetErrorHandler(f func(*Context, error)) {
	s.errorHandler = f
}

func (s *Server) Run() error {
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

func (s *Server) MustRun() {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			logger.Fatal(err)
		}
		s.listener = listener
		logger.Infof("[http] server listening on: %s", s.listener.Addr().String())
	}
	if s.tlsConf != nil {
		if err := s.ServeTLS(s.listener, "", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
		return
	}
	if err := s.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
}

// Start the HTTP server.
func (s *Server) Start(ctx context.Context) {
	go s.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown error: %v", err)
	}

	logger.Info("Server gracefully stopped")
}

// Close immediately closes all active net.
func (s *Server) Close() error {
	logger.Info("[HTTP] server is closing")
	return s.Server.Close()
}

// Shutdown gracefully shuts down the server without interrupting any
// // active connections..
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("[HTTP] server is stopping")
	return s.Server.Shutdown(ctx)
}
