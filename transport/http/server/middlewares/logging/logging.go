package logging

import (
	"time"

	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/transport/http/server"
)

func Middleware() server.Middleware {

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			startTime := time.Now()
			err := next(ctx)
			codeField := logger.Int("status_code", ctx.Response.StatusCode())
			if ctx.Response.StatusCode() < 300 {
				codeField.Background(logger.Green)
			} else if ctx.Response.StatusCode() < 400 {
				codeField.Background(logger.Yellow)
			} else {
				codeField.Background(logger.Red)
			}
			builder := []logger.FieldBuilder{
				logger.String("method", ctx.Request.Method),
				codeField,
				logger.String("path", ctx.Request.URL.Path),
				logger.String("query", ctx.Request.URL.RawQuery),
				logger.String("duration", time.Since(startTime).String()),
			}
			if err != nil {
				builder = append(builder, logger.F("error", err.Error()))
			}

			ctx.Logger().Infos("", builder...)
			return nil
		}
	}
}
