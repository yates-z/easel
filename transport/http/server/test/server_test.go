package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/yates-z/easel/transport/http/server"
)

func Hello(ctx *server.Context) interface{} {
	fmt.Println("Hello World!")
	return "hello world!"
}

func TestServer(t *testing.T) {
	s := server.NewServer(server.Address(":8000"))
	// s.GET("/hello", Hello)
	s.MustRun(context.Background())
}
