package realip

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func TestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 前置处理：在处理 RPC 之前执行
	log.Printf("Unary interceptor: %s", info.FullMethod)

	// 调用处理器函数以处理请求
	resp, err := handler(ctx, req)

	// 后置处理：在处理 RPC 之后执行
	if err != nil {
		log.Printf("Unary interceptor error: %v", err)
	} else {
		log.Printf("Unary interceptor response: %v", resp)
	}

	return resp, err
}
