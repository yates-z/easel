package config

import (
	"sync"

	"github.com/yates-z/easel/core/variant"
	"github.com/yates-z/easel/logger"
)

// Option is config option.
type Option func(*config)

// WithSource with config source.
func WithSource(s Source) Option {
	return func(c *config) {
		c.sources = append(c.sources, s)
	}
}

type config struct {
	mu      sync.RWMutex
	content *Content
	sources []Source
}

// New creates a new configuration instance with the provided options
func New(opts ...Option) *config {
	c := &config{
		content: NewContent(),
	}
	for _, opt := range opts {
		opt(c)
	}

	for _, s := range c.sources {
		content, err := s.Load()
		if err != nil {
			logger.Fatalf(err.Error())
		}
		c.content.Merge(content)
	}
	return c
}

func (c *config) Get(path string) (variant.Variant, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.content.Get(path)
}

// SetInt sets an int value at the specified path.
func (c *config) SetInt(path string, value int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}

// SetFloat64 sets a float64 value at the specified path.
func (c *config) SetFloat64(path string, value float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}

// SetString sets a string value at the specified path.
func (c *config) SetString(path string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}

// SetBool sets a bool value at the specified path.
func (c *config) SetBool(path string, value bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}
