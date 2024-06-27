package logger

import (
	"context"
	"github.com/yates-z/easel/logger/backend"
)

var DefaultLogger Loggable = nil

func Debug(msg ...interface{}) {
	DefaultLogger.Log(DebugLevel, msg...)
}

func Info(msg ...interface{}) {
	DefaultLogger.Log(InfoLevel, msg...)
}

func Warn(msg ...interface{}) {
	DefaultLogger.Log(WarnLevel, msg...)
}

func Error(msg ...interface{}) {
	DefaultLogger.Log(ErrorLevel, msg...)
}

func Fatal(msg ...interface{}) {
	DefaultLogger.Log(FatalLevel, msg...)
}

func Panic(msg ...interface{}) {
	DefaultLogger.Log(PanicLevel, msg...)
}

func Debugf(format string, fmtArgs ...interface{}) {
	DefaultLogger.Logf(DebugLevel, format, fmtArgs...)
}

func Infof(format string, fmtArgs ...interface{}) {
	DefaultLogger.Logf(InfoLevel, format, fmtArgs...)
}

func Warnf(format string, fmtArgs ...interface{}) {
	DefaultLogger.Logf(WarnLevel, format, fmtArgs...)
}

func Errorf(format string, fmtArgs ...interface{}) {
	DefaultLogger.Logf(ErrorLevel, format, fmtArgs...)
}

func Fatalf(format string, fmtArgs ...interface{}) {
	DefaultLogger.Logf(FatalLevel, format, fmtArgs...)
}

func Panicf(format string, fmtArgs ...interface{}) {
	DefaultLogger.Logf(PanicLevel, format, fmtArgs...)
}

func Context(ctx context.Context) Logger {
	return DefaultLogger.(Logger).Context(ctx)
}

func UseDefault() {
	DefaultLogger = NewLogger(
		WithLevel(DebugLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(AnyLevel, " "),
		WithFields(AnyLevel,
			DatetimeField("2006-01-02 15:04:03").Key("datetime").Build(),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Red).Build(),
		),
		WithFields(AnyLevel^ErrorLevel^FatalLevel^PanicLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Build(),
		),
		WithFields(AnyLevel,
			ShortCallerField(true).Key("file").Build(),
			MessageField().Key("msg").Build(),
		),
		WithEncoders(AnyLevel, PlainEncoder),
	)
}

func Use(logger Loggable) {
	DefaultLogger = logger
}
