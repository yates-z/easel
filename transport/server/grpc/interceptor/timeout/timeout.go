package timeout

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// UnaryServerInterceptor returns a new unary client interceptor that sets a timeout on the request context.
func UnaryServerInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		var resp interface{}
		var err error
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		h := func(c context.Context, r interface{}) {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			resp, err = handler(c, r)
			finish <- struct{}{}
		}

		go h(ctx, req)
		select {
		case p := <-panicChan:
			panic(p)
		case _ = <-finish:
			return resp, err
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "request %s timed out.", info.FullMethod)
		}
	}
}
