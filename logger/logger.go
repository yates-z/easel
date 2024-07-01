package logger

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/internal/pool"
	"github.com/yates-z/easel/logger/backend"
	"os"
)

func init() {
	UseDefault()
}

type Loggable interface {
	Log(level LogLevel, args ...any)
	Logf(level LogLevel, format string, fmtArgs ...any)
	Logs(level LogLevel, msg string, fields ...FieldBuilder)
}

type ExtendLoggable interface {
	Debug(args ...any)
	Debugf(format string, fmtArgs ...interface{})
	Debugs(msg string, fields ...FieldBuilder)
	Warn(args ...any)
	Warnf(format string, fmtArgs ...interface{})
	Warns(msg string, fields ...FieldBuilder)
	Info(args ...any)
	Infof(format string, fmtArgs ...interface{})
	Infos(msg string, fields ...FieldBuilder)
	Error(args ...any)
	Errorf(format string, fmtArgs ...interface{})
	Errors(msg string, fields ...FieldBuilder)
	Fatal(args ...any)
	Fatalf(format string, fmtArgs ...interface{})
	Fatals(msg string, fields ...FieldBuilder)
	Panic(args ...any)
	Panicf(format string, fmtArgs ...interface{})
	Panics(msg string, fields ...FieldBuilder)
}

type Logger interface {
	Loggable
	ExtendLoggable
	Backends() []backend.Backend
	Level() LogLevel
	Context(ctx context.Context) Logger
}

type logger struct {
	ctx        context.Context
	level      LogLevel
	entities   map[LogLevel]*logEntity
	entityPool *pool.Pool[*logEntity]
}

func (l *logger) Log(level LogLevel, args ...any) {
	if !l.level.Enabled(level) {
		return
	}
	var msg string
	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			msg = str
		}
	} else {
		msg = fmt.Sprint(args...)
	}

	entity := l.entities[level]

	errs := entity.log(msg)

	l.entities[ErrorLevel].handleError(errs)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(msg)
	}
}

func (l *logger) Logf(level LogLevel, format string, fmtArgs ...any) {
	if !l.level.Enabled(level) {
		return
	}

	if len(fmtArgs) != 0 && format != "" {
		format = fmt.Sprintf(format, fmtArgs...)
	} else if len(fmtArgs) != 0 && format == "" {
		format = fmt.Sprint(fmtArgs...)
	}

	entity := l.entities[level]
	errs := entity.log(format)
	l.entities[ErrorLevel].handleError(errs)

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(format)
	}
}

func (l *logger) Logs(level LogLevel, msg string, fields ...FieldBuilder) {
	if len(fields) == 0 {
		l.Log(level, msg)
		return
	}
	if !l.level.Enabled(level) {
		return
	}

	entity := l.entityPool.Get()
	entity.copy(l.entities[level])
	for _, field := range fields {
		entity.opts.fields = append(entity.opts.fields, field.Build())
	}

	errs := entity.log(msg)
	l.entities[ErrorLevel].handleError(errs)

	for _, field := range fields {
		field.Build().Free()
	}
	entity.clear()
	l.entityPool.Put(entity)

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
	backendMap := map[string]backend.Backend{}
	var backends []backend.Backend
	for _, entity := range l.entities {
		for name, b := range entity.opts.backends {
			if _, ok := backendMap[name]; ok {
				continue
			}
			backendMap[name] = b
			backends = append(backends, b)
		}
	}
	return backends
}

func (l *logger) Level() LogLevel {
	return l.level
}

func (l *logger) Debug(args ...any) {
	l.Log(DebugLevel, args...)
}

func (l *logger) Debugf(format string, fmtArgs ...interface{}) {
	l.Logf(DebugLevel, format, fmtArgs...)
}

func (l *logger) Debugs(msg string, fields ...FieldBuilder) {
	l.Logs(DebugLevel, msg, fields...)
}

func (l *logger) Info(args ...any) {
	l.Log(InfoLevel, args...)
}

func (l *logger) Infof(format string, fmtArgs ...interface{}) {
	l.Logf(InfoLevel, format, fmtArgs...)
}

func (l *logger) Infos(msg string, fields ...FieldBuilder) {
	l.Logs(InfoLevel, msg, fields...)
}

func (l *logger) Warn(args ...any) {
	l.Log(WarnLevel, args...)
}

func (l *logger) Warnf(format string, fmtArgs ...interface{}) {
	l.Logf(WarnLevel, format, fmtArgs...)
}

func (l *logger) Warns(msg string, fields ...FieldBuilder) {
	l.Logs(WarnLevel, msg, fields...)
}

func (l *logger) Error(args ...any) {
	l.Log(ErrorLevel, args...)
}

func (l *logger) Errorf(format string, fmtArgs ...interface{}) {
	l.Logf(ErrorLevel, format, fmtArgs...)
}

func (l *logger) Errors(msg string, fields ...FieldBuilder) {
	l.Logs(ErrorLevel, msg, fields...)
}

func (l *logger) Fatal(args ...any) {
	l.Log(FatalLevel, args...)
}

func (l *logger) Fatalf(format string, fmtArgs ...interface{}) {
	l.Logf(FatalLevel, format, fmtArgs...)
}

func (l *logger) Fatals(msg string, fields ...FieldBuilder) {
	l.Logs(FatalLevel, msg, fields...)
}

func (l *logger) Panic(args ...any) {
	l.Log(PanicLevel, args...)
}

func (l *logger) Panicf(format string, fmtArgs ...interface{}) {
	l.Logf(PanicLevel, format, fmtArgs...)
}

func (l *logger) Panics(msg string, fields ...FieldBuilder) {
	l.Logs(PanicLevel, msg, fields...)
}

func NewLogger(opts ...Option) *logger {
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
		entityPool: pool.New(func() *logEntity {
			return &logEntity{
				level: o.level,
				opts: &entityOptions{
					backends: map[string]backend.Backend{},
				},
			}
		}),
	}
	for _, level := range o.level.EnumIncremental() {
		entity := &logEntity{
			level: level,
			opts:  o.entityOptions[level],
		}
		for _, b := range o.entityOptions[level].backends {
			entity.backends = append(entity.backends, b)
		}
		inst.entities[level] = entity
	}
	return inst
}
