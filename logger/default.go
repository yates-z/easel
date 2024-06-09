package logger

import "context"

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
	return DefaultLogger.Context(ctx)
}

func UseDefault() {
	DefaultLogger = NewLogger()
}

func Use(logger Loggable) {
	DefaultLogger = logger
}
