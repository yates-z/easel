package recovery

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/transport/http/server"
	"runtime"
)

type Option func(*options)

type options struct {
	handler func(*server.Context, any) error
}

// WithHandler with recovery handler.
func WithHandler(h func(*server.Context, any) error) Option {
	return func(o *options) {
		o.handler = h
	}
}

func Middleware(opts ...Option) server.Middleware {
	op := options{
		handler: func(ctx *server.Context, err any) error {
			return errors.New(fmt.Sprintf("unknown request error: %v", err))
		},
	}
	for _, opt := range opts {
		opt(&op)
	}
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := getStack(4)
					logger.Context(ctx).Errorf("%v: %+v\n%s\n", r, ctx.Request, stack)
					err = op.handler(ctx, r)
				}
			}()

			return next(ctx)
		}
	}
}

func getStack(skip int) []byte {
	buf := new(bytes.Buffer)

	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function)
		if !more {
			break
		}
	}
	return buf.Bytes()
}
