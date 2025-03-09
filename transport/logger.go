package transport

import (
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/logger/backend"
)

var Logger = logger.NewLogger(
	logger.WithLevel(logger.DebugLevel),
	logger.WithBackends(logger.AnyLevel, backend.OSBackend().Build()),
	logger.WithSeparator(logger.AnyLevel, "    "),
	logger.WithFields(logger.AnyLevel,
		logger.DatetimeField("2006/01/02 15:04:03").Key("datetime"),
	),
	logger.WithFields(logger.DebugLevel|logger.InfoLevel,
		logger.LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(logger.Green),
	),
	logger.WithFields(logger.WarnLevel,
		logger.LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(logger.Yellow),
	),
	logger.WithFields(logger.ErrorLevel|logger.FatalLevel|logger.PanicLevel,
		logger.LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(logger.Red),
	),
	logger.WithFields(logger.AnyLevel, logger.MessageField().Key("msg")),
	logger.WithEncoders(logger.AnyLevel, logger.PlainEncoder),
)
