package logger

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/logger/backend"
	"os"
)

func init() {
	UseDefault()
}

type Loggable interface {
	Log(level LogLevel, msg ...interface{})
	Logf(level LogLevel, format string, fmtArgs ...interface{})
	Logw(level LogLevel, msg string, fields ...LogField)
}

type Logger interface {
	Loggable
	Backends() []backend.Backend
	Level() LogLevel
	Context(ctx context.Context) Logger
	Debug(msg ...interface{})
	Debugf(format string, fmtArgs ...interface{})
	Debugw(msg string, fields ...LogField)
	Warn(msg ...interface{})
	Warnf(format string, fmtArgs ...interface{})
	Warnw(msg string, fields ...LogField)
	Info(msg ...interface{})
	Infof(format string, fmtArgs ...interface{})
	Infow(msg string, fields ...LogField)
	Error(msg ...interface{})
	Errorf(format string, fmtArgs ...interface{})
	Errorw(msg string, fields ...LogField)
	Fatal(msg ...interface{})
	Fatalf(format string, fmtArgs ...interface{})
	Fatalw(msg string, fields ...LogField)
	Panic(msg ...interface{})
	Panicf(format string, fmtArgs ...interface{})
	Panicw(msg string, fields ...LogField)
}

type logger struct {
	ctx      context.Context
	level    LogLevel
	entities map[LogLevel]*logEntity
}

func (l *logger) Log(level LogLevel, msg ...interface{}) {
	if !l.level.Enabled(level) {
		return
	}
	entity := l.entities[level]
	errs, available := entity.log(fmt.Sprint(msg...))
	l.entities[ErrorLevel].handleError(errs, available)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(fmt.Sprint(msg...))
	}
}

func (l *logger) Logf(level LogLevel, format string, fmtArgs ...interface{}) {
	if !l.level.Enabled(level) {
		return
	}

	if len(fmtArgs) != 0 && format != "" {
		format = fmt.Sprintf(format, fmtArgs...)
	}

	entity := l.entities[level]
	errs, available := entity.log(format)
	l.entities[ErrorLevel].handleError(errs, available)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(format)
	}
}

func (l *logger) Logw(level LogLevel, msg string, fields ...LogField) {
	//TODO implement me
	panic("implement me")
}

func (l *logger) Context(ctx context.Context) Logger {
	l.ctx = ctx
	return l
}

func (l *logger) Backends() []backend.Backend {
	//TODO implement me
	panic("implement me")
}

func (l *logger) Level() LogLevel {
	return l.level
}

func (l *logger) Debug(msg ...interface{}) {
	l.Log(DebugLevel, msg...)
}

func (l *logger) Debugf(format string, fmtArgs ...interface{}) {
	l.Logf(DebugLevel, format, fmtArgs...)
}

func (l *logger) Debugw(msg string, fields ...LogField) {
	l.Logw(DebugLevel, msg, fields...)
}

func (l *logger) Info(msg ...interface{}) {
	l.Log(InfoLevel, msg...)
}

func (l *logger) Infof(format string, fmtArgs ...interface{}) {
	l.Logf(InfoLevel, format, fmtArgs...)
}

func (l *logger) Infow(msg string, fields ...LogField) {
	l.Logw(InfoLevel, msg, fields...)
}

func (l *logger) Warn(msg ...interface{}) {
	l.Log(WarnLevel, msg...)
}

func (l *logger) Warnf(format string, fmtArgs ...interface{}) {
	l.Logf(WarnLevel, format, fmtArgs...)
}

func (l *logger) Warnw(msg string, fields ...LogField) {
	l.Logw(WarnLevel, msg, fields...)
}

func (l *logger) Error(msg ...interface{}) {
	l.Log(ErrorLevel, msg...)
}

func (l *logger) Errorf(format string, fmtArgs ...interface{}) {
	l.Logf(ErrorLevel, format, fmtArgs...)
}

func (l *logger) Errorw(msg string, fields ...LogField) {
	l.Logw(ErrorLevel, msg, fields...)
}

func (l *logger) Fatal(msg ...interface{}) {
	l.Log(FatalLevel, msg...)
}

func (l *logger) Fatalf(format string, fmtArgs ...interface{}) {
	l.Logf(FatalLevel, format, fmtArgs...)
}

func (l *logger) Fatalw(msg string, fields ...LogField) {
	l.Logw(FatalLevel, msg, fields...)
}

func (l *logger) Panic(msg ...interface{}) {
	l.Log(PanicLevel, msg...)
}

func (l *logger) Panicf(format string, fmtArgs ...interface{}) {
	l.Logf(PanicLevel, format, fmtArgs...)
}

func (l *logger) Panicw(msg string, fields ...LogField) {
	l.Logw(PanicLevel, msg, fields...)
}

func NewLogger(opts ...Option) Logger {
	// handle options.
	o := newOptions(InfoLevel)
	for _, opt := range opts {
		opt(o)
	}
	// new logger instance.
	inst := &logger{
		ctx:      context.Background(),
		level:    o.level,
		entities: map[LogLevel]*logEntity{},
	}
	for _, level := range o.level.EnumIncremental() {
		entity := &logEntity{
			level: level,
			opts:  o.entityOptions[level],
		}
		inst.entities[level] = entity
	}
	return inst
}
