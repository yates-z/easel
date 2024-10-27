package recovery

import (
	"bytes"
	"context"
	"fmt"
	"github.com/yates-z/easel/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

type RecoveryHandlerFunc func(context.Context, any) error

type Option func(*options)

type options struct {
	recoveryHandler RecoveryHandlerFunc
}

// WithHandler with recovery handler.
func WithHandler(h RecoveryHandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandler = h
	}
}

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	op := options{
		recoveryHandler: func(ctx context.Context, req interface{}) error {
			return status.Errorf(codes.Unknown, "unknown request error: %v", req)
		},
	}
	for _, opt := range opts {
		opt(&op)
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := getStack(4)
				logger.Context(ctx).Errorf("%v: %+v\n%s\n", r, req, stack)
				err = op.recoveryHandler(ctx, r)
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	op := options{
		recoveryHandler: func(ctx context.Context, req interface{}) error {
			return status.Errorf(codes.Unknown, "unknown request error: %v", req)
		},
	}
	for _, opt := range opts {
		opt(&op)
	}
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 64<<10)
				n := runtime.Stack(stack, false)
				stack = stack[:n]
				logger.Context(stream.Context()).Errorf("%v: \n%s\n", r, stack)
				err = op.recoveryHandler(stream.Context(), r)
			}
		}()

		return handler(srv, stream)
	}
}

func getStack(skip int) []byte {
	buf := new(bytes.Buffer)

	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function)
		if !more {
			break
		}
	}
	return buf.Bytes()
}
