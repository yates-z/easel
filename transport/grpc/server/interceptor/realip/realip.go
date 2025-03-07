package realip

import (
	"context"

	"github.com/yates-z/easel/logger"
	"google.golang.org/grpc"
)

func TestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 前置处理：在处理 RPC 之前执行
	logger.Infof("Unary interceptor: %s", info.FullMethod)

	// 调用处理器函数以处理请求
	resp, err := handler(ctx, req)

	// 后置处理：在处理 RPC 之后执行
	if err != nil {
		logger.Infof("Unary interceptor error: %v", err)
	} else {
		logger.Infof("Unary interceptor response: %v", resp)
	}

	return resp, err
}
