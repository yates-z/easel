package logger

import (
	"easel/logger/backend"
	"fmt"
	"os"
)

type Loggable interface {
	Log(level LogLevel, msg ...interface{})
	Logf(level LogLevel, format string, fmtArgs ...interface{})
}

type Logger interface {
	Loggable
	Backends() []backend.Backend
	Mode() LogMode
	Level() LogLevel
}

type logEntity struct {
	level     LogLevel
	fields    []LogField
	encoders  []Encoder
	separator string
	skipLines int
}

type logger struct {
	opts     Options
	entities map[LogLevel]*logEntity
}

func (l *logger) Level() LogLevel {
	return l.opts.Level
}

func (l *logger) Backends() []backend.Backend {
	return l.opts.Backends
}

func (l *logger) Mode() LogMode {
	//TODO implement me
	panic("implement me")
}

func (l *logger) Log(level LogLevel, msg ...interface{}) {
	if !l.opts.Level.Enabled(level) {
		return
	}
	entity := l.entities[level]

	var logs, plogs []byte
	for _, encoder := range entity.encoders {
		log, plog := encoder.Encode(fmt.Sprint(msg...))
		logs = append(logs, []byte(log)...)
		plogs = append(plogs, []byte(plog)...)
	}

	var errs string
	var available []backend.Backend
	for _, b := range l.opts.Backends {
		var err error
		if b.AllowANSI() {
			_, err = b.Write(plogs)
		} else {
			_, err = b.Write(logs)
		}
		if err != nil {
			errs += fmt.Sprintf("%T writer error: %s.", b, err)
		} else {
			available = append(available, b)
		}
	}
	if len(errs) != 0 {
		// broadcast errors
		entity = l.entities[ErrorLevel]

		for _, encoder := range entity.encoders {
			log, plog := encoder.Encode(errs)
			logs = append(logs, []byte(log)...)
			plogs = append(plogs, []byte(plog)...)
		}

		for _, b := range available {
			if b.AllowANSI() {
				_, _ = b.Write(plogs)
			} else {
				_, _ = b.Write(logs)
			}
		}
	}

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(fmt.Sprint(msg...))
	}

}

func (l *logger) Logf(level LogLevel, format string, fmtArgs ...interface{}) {
	if !l.opts.Level.Enabled(level) {
		return
	}

	if len(fmtArgs) != 0 && format != "" {
		format = fmt.Sprintf(format, fmtArgs...)
	}

	entity := l.entities[level]
	var logs, plogs []byte
	for _, encoder := range entity.encoders {
		log, plog := encoder.Encode(format)
		logs = append(logs, []byte(log)...)
		plogs = append(plogs, []byte(plog)...)
	}

	var errs string
	var available []backend.Backend
	for _, b := range l.opts.Backends {
		var err error
		if b.AllowANSI() {
			_, err = b.Write(plogs)
		} else {
			_, err = b.Write(logs)
		}
		if err != nil {
			errs += fmt.Sprintf("%T writer error: %s.", b, err)
		} else {
			available = append(available, b)
		}
	}
	if len(errs) != 0 {
		// broadcast errors
		entity = l.entities[ErrorLevel]
		for _, encoder := range entity.encoders {
			log, plog := encoder.Encode(errs)
			logs = append(logs, []byte(log)...)
			plogs = append(plogs, []byte(plog)...)
		}

		for _, b := range available {
			if b.AllowANSI() {
				_, _ = b.Write(plogs)
			} else {
				_, _ = b.Write(logs)
			}
		}
	}

	if level.Eq(FatalLevel) {
		os.Exit(1)
	}
	if level.Eq(PanicLevel) {
		panic(format)
	}
}

func NewLogger(opts ...Option) Logger {
	// handle options.
	options := Options{
		Level:     InfoLevel,
		Separator: " ",
		SkipLines: 0,
	}
	for _, opt := range opts {
		opt(&options)
	}
	// fill attributes of options by mode.
	if len(options.Backends) == 0 {
		options.Backends = append(options.Backends, backend.OSBackend().Build())
	}
	if len(options.Fields) == 0 {
		options.Fields = append(
			options.Fields,
			LevelField("level").Upper(true).Build(),
			DatetimeField("datetime").Build(),
			ShortFileField("file").Build(),
			BodyField("body").Build(),
		)
	}
	if len(options.Encoders) == 0 {
		options.Encoders = append(options.Encoders, PlainEncoder())
	}
	// new logger instance.
	inst := &logger{
		opts:     options,
		entities: map[LogLevel]*logEntity{},
	}
	for _, level := range options.Level.Enum() {
		entity := &logEntity{
			level:     level,
			separator: options.Separator,
			skipLines: options.SkipLines,
		}
		inst.entities[level] = entity

		for _, constructor := range options.Fields {
			entity.fields = append(entity.fields, constructor(entity))
		}
		for _, constructor := range options.Encoders {
			entity.encoders = append(entity.encoders, constructor(entity))
		}
	}
	return inst
}
