package logger

import "github.com/yates-z/easel/logger/backend"

type Option func(*options)

type options struct {
	level         LogLevel
	entityOptions map[LogLevel]*entityOptions
}

func newOptions(defaultLevel LogLevel) *options {
	o := &options{
		level:         defaultLevel,
		entityOptions: map[LogLevel]*entityOptions{},
	}
	for _, level := range DebugLevel.EnumIncremental() {
		o.entityOptions[level] = &entityOptions{
			backends: map[backend.Backend]struct{}{},
		}
	}
	return o
}

type entityOptions struct {
	separator string
	skipLines int
	fields    []LogField
	encoders  []Encoder
	backends  map[backend.Backend]struct{}
}

// WithLevel set default level for the logger.
func WithLevel(level LogLevel) Option {
	return func(opts *options) {
		opts.level = level
	}
}

// WithSeparator only have effection when a specific encoder
// such as PlainEncoder is used.
func WithSeparator(collection LogLevel, separator string) Option {
	return func(opts *options) {
		for _, level := range collection.Enum() {
			opts.entityOptions[level].separator = separator
		}
	}
}

// WithSkipLines .
func WithSkipLines(collection LogLevel, c int) Option {
	return func(opts *options) {
		for _, level := range collection.Enum() {
			opts.entityOptions[level].skipLines = c
		}
	}
}

func WithFields(collection LogLevel, field ...LogField) Option {
	return func(opts *options) {
		for _, level := range collection.Enum() {
			opts.entityOptions[level].fields = append(opts.entityOptions[level].fields, field...)
		}
	}
}

func WithEncoders(collection LogLevel, encoders ...Encoder) Option {
	return func(opts *options) {
		for _, level := range collection.Enum() {
			opts.entityOptions[level].encoders = append(opts.entityOptions[level].encoders, encoders...)
		}
	}
}

// WithBackends set default backend for the logger.
func WithBackends(collection LogLevel, backends ...backend.Backend) Option {
	return func(opts *options) {
		for _, level := range collection.Enum() {
			for _, b := range backends {
				opts.entityOptions[level].backends[b] = struct{}{}
			}
		}
	}
}
