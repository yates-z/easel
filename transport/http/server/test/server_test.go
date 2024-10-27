package test

import (
	"context"
	"github.com/yates-z/easel/transport/grpc/server/test/api"
	"github.com/yates-z/easel/transport/http/server"
	"github.com/yates-z/easel/transport/http/server/adapter"
	"github.com/yates-z/easel/transport/http/server/middlewares/recovery"
	"net/http"
	"testing"
)

type HelloService struct {
}

func (s *HelloService) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {

	return &api.HelloResponse{Replay: "hello, " + in.Name}, nil
}

func Hello(ctx *server.Context) error {
	panic("pppanic")
	return ctx.JSON(http.StatusOK, map[string]string{"hello": ctx.Param("name")})
}

func TestServer(t *testing.T) {
	s := server.New(
		server.Address(":8000"),
		server.ShowInfo(),
		server.Middlewares(recovery.Middleware()),
	)
	service := &HelloService{}
	s.GET("/hello/{name}", Hello)
	s.POST("/hello/{$}", adapter.GRPC(service.SayHello))
	s.MustRun()
}
