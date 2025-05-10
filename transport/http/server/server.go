package server

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/yates-z/easel/core/pool"
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/transport"
	templ "github.com/yates-z/easel/transport/http/server/template"
	"github.com/yates-z/easel/utils/host"
	"html/template"
	"net"
	"net/http"
	"net/url"
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

	log          logger.Logger
	ctxPool      *pool.Pool[*Context]
	middlewares  []Middleware
	showInfo     bool
	htmlTempl    *templ.HTMLTemplate
	errorHandler func(ctx *Context, err error)
}

func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		network:   "tcp",
		address:   ":80",
		log:       transport.Logger,
		showInfo:  false,
		htmlTempl: templ.New(),
	}
	server.Router = NewRouter(server)
	server.ctxPool = pool.New(func() *Context {
		return newContext(server)
	})
	server.errorHandler = func(ctx *Context, err error) {
		if ctx.GetHeader("Content-Type") == "application/json" {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}
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

func (s *Server) Start(ctx context.Context) error {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.listener = listener
	}
	// http.Server.BaseContext：定义服务器级别的上下文，所有请求都会继承该上下文。
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	s.log.Infof("[http] server listening on: %s", s.listener.Addr().String())
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.listener, "", "")
	} else {
		err = s.Serve(s.listener)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) MustStart(ctx context.Context) {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			s.log.Fatal(err)
		}
		s.listener = listener
		s.log.Infof("[http] server listening on: %s", s.listener.Addr().String())
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	if s.tlsConf != nil {
		if err := s.ServeTLS(s.listener, "", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal(err)
		}
		return
	}
	if err := s.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatal(err)
	}
}

func (s *Server) Endpoint() (*url.URL, error) {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			return nil, err
		}
		s.listener = listener
	}
	addr, err := host.Extract(s.address, s.listener)
	if err != nil {
		return nil, err
	}
	scheme := "http"
	if s.tlsConf != nil {
		scheme = "https"
	}
	return &url.URL{Scheme: scheme, Host: addr}, nil
}

// Close immediately closes all active net.
func (s *Server) Close() error {
	s.log.Info("[HTTP] server is closing")
	return s.Server.Close()
}

// Stop gracefully shuts down the server without interrupting any
// active connections.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[HTTP] server is stopping")
	err := s.Server.Shutdown(ctx)
	if err != nil {
		if ctx.Err() != nil {
			s.log.Warn("[HTTP] server couldn't stop gracefully in time, doing force stop")
			err = s.Close()
		}
	}
	return err
}
