package client

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"slices"
)

type DialOption func(opt *dialOptions)

func WithTarget(target string) DialOption {
	return func(o *dialOptions) {
		o.target = target
	}
}

func WithInsecure() DialOption {
	return func(o *dialOptions) {
		o._opts = append(o._opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) DialOption {
	return func(o *dialOptions) {
		o._opts = append(o._opts, grpc.WithTransportCredentials(credentials.NewTLS(c)))
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the client.
func UnaryInterceptor(in ...grpc.UnaryClientInterceptor) DialOption {
	return func(o *dialOptions) {
		o.unaryInterceptors = in
	}
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the client.
func StreamInterceptor(in ...grpc.StreamClientInterceptor) DialOption {
	return func(o *dialOptions) {
		o.streamInterceptors = in
	}
}

// GRPCOptions with grpc options.
func GRPCOptions(opts ...grpc.DialOption) DialOption {
	return func(o *dialOptions) {
		o._opts = slices.Concat(o._opts, opts)
	}
}

type dialOptions struct {
	target             string
	unaryInterceptors  []grpc.UnaryClientInterceptor
	streamInterceptors []grpc.StreamClientInterceptor
	_opts              []grpc.DialOption
}

func NewClient(opts ...DialOption) (*grpc.ClientConn, error) {
	options := &dialOptions{}
	for _, o := range opts {
		o(options)
	}
	options._opts = append(options._opts, grpc.WithChainUnaryInterceptor(options.unaryInterceptors...))
	options._opts = append(options._opts, grpc.WithChainStreamInterceptor(options.streamInterceptors...))

	return grpc.NewClient(options.target, options._opts...)
}
