package test

import (
	"context"
	"github.com/yates-z/easel/transport/grpc/server"
	"github.com/yates-z/easel/transport/grpc/server/compressor/zlib"
	"github.com/yates-z/easel/transport/grpc/server/interceptor/ratelimit"
	"github.com/yates-z/easel/transport/grpc/server/interceptor/recovery"
	"github.com/yates-z/easel/transport/grpc/server/interceptor/timeout"
	"github.com/yates-z/easel/transport/grpc/server/test/api"
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
	panic("implement me")
	return &api.HelloResponse{Replay: "hello"}, nil
}

func TestServeGrpc(t *testing.T) {
	s := server.NewServer(
		server.Address("0.0.0.0:9100"),
		server.Compressor(zlib.New()),
		server.UnaryInterceptor(
			//ratelimit.UnaryServerInterceptor(ratelimit.NewTokenBucket(time.Second, 1, 10)),
			ratelimit.UnaryServerInterceptor(ratelimit.NewLeakyBucket(1, 2)),
			timeout.UnaryServerInterceptor(5*time.Second),
			recovery.UnaryServerInterceptor(),
		),
		server.StreamInterceptor(recovery.StreamServerInterceptor()),
		server.AllowReflection(true),
	)
	api.RegisterGreeterServer(s, &Server2{})
	s.MustRun()

	//s := grpc.NewServer()
	//api.RegisterGreeterServer(s, &Server2{})
	//lis, _ := net.Listen("tcp", ":8888")
	//_ = s.Serve(lis)
}
