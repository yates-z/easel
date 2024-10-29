package logging

import (
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/logger/backend"
	"github.com/yates-z/easel/transport/http/server"
	"strconv"
	"time"
)

func Middleware() server.Middleware {
	log := logger.NewLogger(
		logger.WithLevel(logger.DebugLevel),
		logger.WithBackends(logger.AnyLevel, backend.OSBackend().Build()),
		logger.WithSeparator(logger.AnyLevel, "    "),
		logger.WithFields(logger.AnyLevel,
			logger.DatetimeField("2006/01/02 15:04:03").Key("datetime"),
			logger.LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(logger.Green),
		),
		logger.WithEncoders(logger.AnyLevel, logger.PlainEncoder),
	)

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			startTime := time.Now()
			err := next(ctx)
			codeField := logger.F("status_code", strconv.Itoa(ctx.Response.StatusCode()))
			if ctx.Response.StatusCode() < 300 {
				codeField.Background(logger.Green)
			} else if ctx.Response.StatusCode() < 400 {
				codeField.Background(logger.Yellow)
			} else {
				codeField.Background(logger.Red)
			}
			builder := []logger.FieldBuilder{
				logger.F("method", ctx.Request.Method),
				codeField,
				logger.F("path", ctx.Request.URL.Path),
				logger.F("query", ctx.Request.URL.RawQuery),
				logger.F("duration", time.Since(startTime).String()),
			}
			if err != nil {
				builder = append(builder, logger.F("error", err.Error()))
			}

			log.Infos("", builder...)
			return nil
		}
	}
}
