package grpc

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"slices"
)

type ServerOption func(*Server)
type EmptyCompressor string

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
		s._opts = append(s._opts, grpc.Creds(credentials.NewTLS(c)))
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptors = in
	}
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func StreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptors = in
	}
}

// AllowReflection determines whether to use reflection.
func AllowReflection(allow bool) ServerOption {
	return func(s *Server) {
		s.allowReflection = allow
	}
}

// AllowHealthCheck determines whether to use health.
func AllowHealthCheck(allow bool) ServerOption {
	return func(s *Server) {
		s.allowHealthCheck = allow
	}
}

// Compressor with server address.
func Compressor(compressor encoding.Compressor) ServerOption {
	return func(s *Server) {
		s.compressor = compressor
	}
}

// GRPCOptions with grpc options.
func GRPCOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s._opts = slices.Concat(s._opts, opts)
	}
}

type Server struct {
	*grpc.Server

	// network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	// Check net.Listen for more detail.
	network string

	// address optionally specifies the address for the server to listen on.
	// in the form "host:port". If empty, ":http" (port 80) is used.
	address  string
	listener net.Listener

	// interceptors collect grpc interceptors that have same effect on both unary and stream req.
	interceptors []grpc.UnaryServerInterceptor

	// unaryInterceptors are hooks to intercept the execution of a unary RPC on the server.
	unaryInterceptors []grpc.UnaryServerInterceptor

	// streamInterceptors are hooks to intercept the execution of a streaming RPC on the server.
	streamInterceptors []grpc.StreamServerInterceptor

	// allowReflection determines whether register reflection service on server.
	allowReflection bool

	// allowHealthCheck determines whether register health service on server.
	allowHealthCheck bool

	// health is a health service.
	health *health.Server

	// compressor will compress grpc message.
	// You can use custom compressor or build-in 'gzip'.
	compressor encoding.Compressor

	//_opts are grpc options for init.
	_opts []grpc.ServerOption
}

func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		network:          "tcp",
		address:          ":0",
		allowReflection:  false,
		allowHealthCheck: true,
		health:           health.NewServer(),
		compressor:       nil,
	}
	for _, opt := range opts {
		opt(server)
	}

	if server.compressor != nil {
		encoding.RegisterCompressor(server.compressor)
	}
	server._opts = append(server._opts, grpc.ChainUnaryInterceptor(server.unaryInterceptors...))
	server._opts = append(server._opts, grpc.ChainStreamInterceptor(server.streamInterceptors...))

	server.Server = grpc.NewServer(server._opts...)

	if server.allowHealthCheck {
		grpc_health_v1.RegisterHealthServer(server, server.health)
	}
	if server.allowReflection {
		reflection.Register(server)
	}
	return server
}

func (s *Server) Listen(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	s.listener = listener
	log.Printf("[gRPC] server listening on: %s", listener.Addr().String())
	return nil
}

func (s *Server) Run() error {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.listener = listener
		log.Printf("[gRPC] server listening on: %s", s.listener.Addr().String())
	}
	s.health.Resume()
	return s.Serve(s.listener)
}

func (s *Server) MustRun() {
	if s.listener == nil {
		listener, err := net.Listen(s.network, s.address)
		if err != nil {
			log.Fatal(err)
		}
		s.listener = listener
		log.Printf("[gRPC] server listening on: %s", s.listener.Addr().String())
	}
	s.health.Resume()
	log.Fatal(s.Serve(s.listener))
}

func (s *Server) Stop() error {
	s.health.Shutdown()
	s.GracefulStop()
	log.Println("[gRPC] server stopping")
	return nil
}
