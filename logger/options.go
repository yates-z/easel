package logger

import (
	"github.com/yates-z/easel/logger/backend"
)

type Option func(*Options)

type Options struct {
	Backends  []backend.Backend
	Level     LogLevel
	Fields    []FieldConstructor
	Separator string
	SkipLines int
	Encoders  []EncoderConstructor
}

// WithLevel set default level for the logger.
func WithLevel(level LogLevel) Option {
	return func(opts *Options) {
		opts.Level = level
	}
}

// WithBackends set default backend for the logger.
func WithBackends(backend ...backend.Backend) Option {
	return func(opts *Options) {
		opts.Backends = append(opts.Backends, backend...)
	}
}

// WithSkipLines .
func WithSkipLines(c int) Option {
	return func(opts *Options) {
		opts.SkipLines = c
	}
}

func WithFields(field ...FieldConstructor) Option {
	return func(opts *Options) {
		opts.Fields = append(opts.Fields, field...)
	}
}

func WithSeparator(Separator string) Option {
	return func(opts *Options) {
		opts.Separator = Separator
	}
}

func WithEncoders(encoders ...EncoderConstructor) Option {
	return func(opts *Options) {
		opts.Encoders = encoders
	}
}
