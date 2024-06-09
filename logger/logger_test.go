package logger

import (
	"github.com/yates-z/easel/logger/backend"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(
		WithLevel(InfoLevel),
		WithSkipLines(1),
		WithBackends(backend.OSBackend().Build(), backend.DefaultFileBackend().Build()),
		WithFields(
			LevelField("level").Upper(true).Build(),
			DatetimeField("datetime").Build(),
			ShortFileField("file").Build(),
			BodyField("body").Build(),
		),
		WithSeparator(" "),
		WithEncoders(PlainEncoder(), JsonEncoder()),
	)
	logger.Log(DebugLevel, "this is a test")
	logger.Log(WarnLevel, "this is a test", 123, 456)
	logger.Logf(ErrorLevel, "this is a %s test", "error")
	logger.Logf(PanicLevel, "fatal..... %d", 666)
}

func TestColor(t *testing.T) {
	logger := NewLogger(
		WithLevel(InfoLevel),
		WithSkipLines(1),
		WithBackends(
			backend.OSBackend().Build(),
			backend.DefaultFileBackend().Filename("log.log").Build(),
		),
		WithFields(
			LevelField("level").Upper(true).Color(Red).Background(Blue).Build(),
			DatetimeField("datetime").Color(Green).Build(),
			ShortFileField("file").Color(Black).Background(Magenta).Build(),
			BodyField("body").Background(Yellow).Build(),
		),
		WithSeparator("   "),
		WithEncoders(JsonEncoder(), PlainEncoder(), LogFmtEncoder()),
	)
	logger.Log(DebugLevel, "this is a test")
	logger.Log(WarnLevel, "this is a test", 123, 456, nil)
	logger.Logf(ErrorLevel, "this is a %s test", "error")
}
