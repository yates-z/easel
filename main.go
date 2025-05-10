package main

import (
	"context"
	"github.com/yates-z/easel/transport/http/server"
	"github.com/yates-z/easel/transport/http/server/middlewares/logging"
	"github.com/yates-z/easel/transport/http/server/middlewares/recovery"
	"net/http"
)

type R struct {
	Code int
	Data string
}

func Hello(ctx *server.Context) error {
	//fmt.Println("Hello World!")
	return ctx.JSON(http.StatusOK, R{200, "Hello World!"})
}

func main() {

	s := server.NewServer(
		server.Address("localhost:8000"),
		server.ShowInfo(),
		server.Middlewares(
			recovery.Middleware(),
			logging.Middleware(),
			//session.Middleware(session2.NewSessionManager(
			//	session2.NewCacheSessionBackend(1, 1000, time.Minute),
			//	session2.WithTTL(30*time.Minute),
			//)),
		),
	)
	s.GET("/hello1", Hello)
	s.Start(context.Background())
}
