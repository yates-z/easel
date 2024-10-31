package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yates-z/easel/transport/http/server"
	"github.com/yates-z/easel/transport/http/server/middlewares/logging"
	"github.com/yates-z/easel/transport/http/server/middlewares/recovery"
)

type R struct {
	Code int
	Data string
}

func Hello(ctx *server.Context) error {
	fmt.Println("Hello World!")
	return ctx.JSON(http.StatusOK, R{200, "Hello World!"})
}

func main() {

	s := server.New(
		server.Address("localhost:8000"),
		server.ShowInfo(),
		server.Middlewares(recovery.Middleware(), logging.Middleware()),
	)
	s.GET("/hello", Hello)
	s.Start(context.Background())
}
