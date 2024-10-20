package main

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/transport/http/server"
	"net/http"
)

type R struct {
	Code int
	Data string
}

func Hello(ctx *server.Context) (interface{}, error) {
	fmt.Println("Hello World!")
	return ctx.JSON(http.StatusOK, R{200, "Hello World!"}), nil
}

func main() {

	s := server.New(server.Address("localhost:8000"), server.ShowInfo(true))
	s.GET("/hello", Hello)
	s.MustRun(context.Background())
}
