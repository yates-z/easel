package client

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/api"
	"github.com/yates-z/easel/transport/grpc/client/interceptor/retry"
	"github.com/yates-z/easel/transport/grpc/client/interceptor/timeout"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewInsecureClient(
		"127.0.0.1:9100",
		UnaryInterceptor(
			timeout.UnaryClientInterceptor(0),
			retry.UnaryClientInterceptor(5, time.Second, retry.WithPerRetryTimeout(time.Second)),
		),
	)

	if err != nil {
		panic(err)
	}

	c := api.NewGreeterClient(client)
	res, err := c.SayHello(context.Background(), &api.HelloRequest{Name: ""})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
