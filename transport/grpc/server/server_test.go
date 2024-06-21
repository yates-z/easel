package server

import (
	"context"
	"github.com/yates-z/easel/api"
	"github.com/yates-z/easel/transport/grpc/server/compressor/zlib"
	"github.com/yates-z/easel/transport/grpc/server/interceptor/ratelimit"
	"github.com/yates-z/easel/transport/grpc/server/interceptor/recovery"
	"github.com/yates-z/easel/transport/grpc/server/interceptor/timeout"
	"testing"
	"time"
)

type Server2 struct {
	api.UnimplementedGreeterServer
}

func (s *Server2) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	// 从上下文中提取 metadata
	//if md, ok := metadata.FromIncomingContext(ctx); ok {
	//	for key, values := range md {
	//		for _, value := range values {
	//			log.Printf("Metadata key: %s, value: %s", key, value)
	//		}
	//	}
	//}
	time.Sleep(3 * time.Second)
	//panic("implement me")
	return &api.HelloResponse{Replay: "hello"}, nil
}

func TestServeGrpc(t *testing.T) {
	server := NewServer(
		Address("0.0.0.0:9100"),
		Compressor(zlib.New()),
		UnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			timeout.UnaryServerInterceptor(5*time.Second),
			//ratelimit.UnaryServerInterceptor(ratelimit.NewTokenBucket(time.Second, 1, 10)),
			ratelimit.UnaryServerInterceptor(ratelimit.NewLeakyBucket(1, 2)),
		),
		StreamInterceptor(recovery.StreamServerInterceptor()),
		AllowReflection(true),
	)
	api.RegisterGreeterServer(server, &Server2{})
	server.MustRun()

	//s := grpc.NewServer()
	//api.RegisterGreeterServer(s, &Server2{})
	//lis, _ := net.Listen("tcp", ":8888")
	//_ = s.Serve(lis)
}
