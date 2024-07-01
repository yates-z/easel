package logger

import (
	"context"
	"github.com/yates-z/easel/logger/backend"
)

var DefaultLogger *logger = nil

func Debug(args ...any) {
	DefaultLogger.Log(DebugLevel, args...)
}

func Info(args ...any) {
	DefaultLogger.Log(InfoLevel, args...)
}

func Warn(args ...any) {
	DefaultLogger.Log(WarnLevel, args...)
}

func Error(args ...any) {
	DefaultLogger.Log(ErrorLevel, args...)
}

func Fatal(args ...any) {
	DefaultLogger.Log(FatalLevel, args...)
}

func Panic(args ...any) {
	DefaultLogger.Log(PanicLevel, args...)
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

func Debugs(msg string, fields ...FieldBuilder) {
	DefaultLogger.Logs(DebugLevel, msg, fields...)
}

func Infos(msg string, fields ...FieldBuilder) {
	DefaultLogger.Logs(InfoLevel, msg, fields...)
}

func Warns(msg string, fields ...FieldBuilder) {
	DefaultLogger.Logs(WarnLevel, msg, fields...)
}

func Errors(msg string, fields ...FieldBuilder) {
	DefaultLogger.Logs(ErrorLevel, msg, fields...)
}

func Fatals(msg string, fields ...FieldBuilder) {
	DefaultLogger.Logs(FatalLevel, msg, fields...)
}

func Panics(msg string, fields ...FieldBuilder) {
	DefaultLogger.Logs(PanicLevel, msg, fields...)
}

func Context(ctx context.Context) Logger {
	return DefaultLogger.Context(ctx)
}

func UseDefault() {
	DefaultLogger = NewLogger(
		WithLevel(DebugLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(AnyLevel, "    "),
		WithFields(AnyLevel,
			DatetimeField("2006/01/02 15:04:03").Key("datetime"),
		),
		WithFields(DebugLevel|InfoLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Green),
		),
		WithFields(WarnLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Yellow),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Red),
		),
		WithFields(AnyLevel,
			MessageField().Key("msg"),
			CallerField(true, true).Key("caller"),
		),
		WithEncoders(AnyLevel, PlainEncoder),
	)
}
