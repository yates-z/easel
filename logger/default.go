package logger

import (
	"context"
	"github.com/yates-z/easel/logger/backend"
)

var DefaultLogger Loggable = nil

func Debug(msg string, args ...any) {
	DefaultLogger.Log(DebugLevel, msg, args...)
}

func Info(msg string, args ...any) {
	DefaultLogger.Log(InfoLevel, msg, args...)
}

func Warn(msg string, args ...any) {
	DefaultLogger.Log(WarnLevel, msg, args...)
}

func Error(msg string, args ...any) {
	DefaultLogger.Log(ErrorLevel, msg, args...)
}

func Fatal(msg string, args ...any) {
	DefaultLogger.Log(FatalLevel, msg, args...)
}

func Panic(msg string, args ...any) {
	DefaultLogger.Log(PanicLevel, msg, args...)
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

func Debugs(msg string, fields ...LogField) {
	DefaultLogger.Logs(DebugLevel, msg, fields...)
}

func Infos(msg string, fields ...LogField) {
	DefaultLogger.Logs(InfoLevel, msg, fields...)
}

func Warns(msg string, fields ...LogField) {
	DefaultLogger.Logs(WarnLevel, msg, fields...)
}

func Errors(msg string, fields ...LogField) {
	DefaultLogger.Logs(ErrorLevel, msg, fields...)
}

func Fatals(msg string, fields ...LogField) {
	DefaultLogger.Logs(FatalLevel, msg, fields...)
}

func Panics(msg string, fields ...LogField) {
	DefaultLogger.Logs(PanicLevel, msg, fields...)
}

func Context(ctx context.Context) Logger {
	return DefaultLogger.(Logger).Context(ctx)
}

func UseDefault() {
	DefaultLogger = NewLogger(
		WithLevel(DebugLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(AnyLevel, "    "),
		WithFields(AnyLevel,
			DatetimeField("2006/01/02 15:04:03").Key("datetime").Build(),
		),
		WithFields(DebugLevel|InfoLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Green).Build(),
		),
		WithFields(WarnLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Yellow).Build(),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Red).Build(),
		),
		WithFields(AnyLevel,
			MessageField().Key("msg").Build(),
			CallerField(true, true).Key("caller").Build(),
		),
		WithEncoders(AnyLevel, PlainEncoder),
	)
}

func Use(logger Loggable) {
	DefaultLogger = logger
}
