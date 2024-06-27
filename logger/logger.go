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
	Log(level LogLevel, msg string, args ...any)
	Logf(level LogLevel, format string, fmtArgs ...any)
	Logs(level LogLevel, msg string, fields ...LogField)
}

type Logger interface {
	Loggable
	Backends() []backend.Backend
	Level() LogLevel
	Context(ctx context.Context) Logger
	Debug(msg string, args ...any)
	Debugf(format string, fmtArgs ...interface{})
	Debugs(msg string, fields ...LogField)
	Warn(msg string, args ...any)
	Warnf(format string, fmtArgs ...interface{})
	Warns(msg string, fields ...LogField)
	Info(msg string, args ...any)
	Infof(format string, fmtArgs ...interface{})
	Infos(msg string, fields ...LogField)
	Error(msg string, args ...any)
	Errorf(format string, fmtArgs ...interface{})
	Errors(msg string, fields ...LogField)
	Fatal(msg string, args ...any)
	Fatalf(format string, fmtArgs ...interface{})
	Fatals(msg string, fields ...LogField)
	Panic(msg string, args ...any)
	Panicf(format string, fmtArgs ...interface{})
	Panics(msg string, fields ...LogField)
}

type logger struct {
	ctx      context.Context
	level    LogLevel
	entities map[LogLevel]*logEntity
}

func (l *logger) Log(level LogLevel, msg string, args ...any) {
	if !l.level.Enabled(level) {
		return
	}
	entity := l.entities[level]
	msg = fmt.Sprint(msg, fmt.Sprint(args...))

	errs, available := entity.log(&msg)
	l.entities[ErrorLevel].handleError(errs, available)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(fmt.Sprint(msg, args))
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
	errs, available := entity.log(&format)
	l.entities[ErrorLevel].handleError(errs, available)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(format)
	}
}

func (l *logger) Logs(level LogLevel, msg string, fields ...LogField) {
	if !l.level.Enabled(level) {
		return
	}
	entity := l.entities[level].copy()
	entity.opts.fields = append(entity.opts.fields, fields...)

	errs, available := entity.log(&msg)
	l.entities[ErrorLevel].handleError(errs, available)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(msg)
	}
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

func (l *logger) Debug(msg string, args ...any) {
	l.Log(DebugLevel, msg, args...)
}

func (l *logger) Debugf(format string, fmtArgs ...interface{}) {
	l.Logf(DebugLevel, format, fmtArgs...)
}

func (l *logger) Debugs(msg string, fields ...LogField) {
	l.Logs(DebugLevel, msg, fields...)
}

func (l *logger) Info(msg string, args ...any) {
	l.Log(InfoLevel, msg, args...)
}

func (l *logger) Infof(format string, fmtArgs ...interface{}) {
	l.Logf(InfoLevel, format, fmtArgs...)
}

func (l *logger) Infos(msg string, fields ...LogField) {
	l.Logs(InfoLevel, msg, fields...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.Log(WarnLevel, msg, args...)
}

func (l *logger) Warnf(format string, fmtArgs ...interface{}) {
	l.Logf(WarnLevel, format, fmtArgs...)
}

func (l *logger) Warns(msg string, fields ...LogField) {
	l.Logs(WarnLevel, msg, fields...)
}

func (l *logger) Error(msg string, args ...any) {
	l.Log(ErrorLevel, msg, args...)
}

func (l *logger) Errorf(format string, fmtArgs ...interface{}) {
	l.Logf(ErrorLevel, format, fmtArgs...)
}

func (l *logger) Errors(msg string, fields ...LogField) {
	l.Logs(ErrorLevel, msg, fields...)
}

func (l *logger) Fatal(msg string, args ...any) {
	l.Log(FatalLevel, msg, args...)
}

func (l *logger) Fatalf(format string, fmtArgs ...interface{}) {
	l.Logf(FatalLevel, format, fmtArgs...)
}

func (l *logger) Fatals(msg string, fields ...LogField) {
	l.Logs(FatalLevel, msg, fields...)
}

func (l *logger) Panic(msg string, args ...any) {
	l.Log(PanicLevel, msg, args...)
}

func (l *logger) Panicf(format string, fmtArgs ...interface{}) {
	l.Logf(PanicLevel, format, fmtArgs...)
}

func (l *logger) Panics(msg string, fields ...LogField) {
	l.Logs(PanicLevel, msg, fields...)
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
