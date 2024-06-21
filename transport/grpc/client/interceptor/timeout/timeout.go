package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// UnaryClientInterceptor returns a new unary client interceptor that sets a timeout on the request context.
func UnaryClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout > 0 {
			timedCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return invoker(timedCtx, method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
