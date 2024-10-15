package server

import (
	"context"
	"net/http"
	"time"
)

type HandlerFunc func(*Context) interface{}

var _ context.Context = (*Context)(nil)

// Context is an HTTP request Context. It defines core functions sets of this http.server.
type Context struct {
	request  *http.Request
	response http.ResponseWriter
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	ctx := c.request.Context()
	return ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	ctx := c.request.Context()
	return ctx.Done()
}

func (c *Context) Err() error {
	ctx := c.request.Context()
	return ctx.Err()
}

func (c *Context) Value(key any) any {
	ctx := c.request.Context()
	return ctx.Value(key)
}

func newContext(req *http.Request, resp http.ResponseWriter) *Context {
	return &Context{
		request:  req,
		response: resp,
	}
}
