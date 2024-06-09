package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Limit(ctx context.Context) error
}

func UnaryServerInterceptor(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := limiter.Limit(ctx); err != nil {
			return nil, status.Errorf(
				codes.ResourceExhausted,
				"%s unvailable due to rate limit exceeded, please retry later. %s", info.FullMethod, err,
			)
		}
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(limiter Limiter) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := limiter.Limit(stream.Context()); err != nil {
			return status.Errorf(
				codes.ResourceExhausted,
				"%s unvailable due to rate limit exceeded, please retry later. %s", info.FullMethod, err,
			)
		}
		return handler(srv, stream)
	}
}
