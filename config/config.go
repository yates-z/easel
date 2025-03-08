package config

import (
	"sync"
	"time"

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

func (c *config) Get(path string) variant.Variant {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.content.Get(path)
}

func (c *config) GetDefault(path string, _default any) variant.Variant {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.content.GetDefault(path, _default)
}

func (c *config) GetBool(path string, _default bool) bool {
	v := c.GetDefault(path, _default)
	return v.ToBool()
}

func (c *config) GetString(path string, _default string) string {
	v := c.GetDefault(path, _default)
	return v.ToString()
}

func (c *config) GetInt(path string, _default int) int {
	v := c.GetDefault(path, _default)
	return v.ToInt()
}

func (c *config) GetUint(path string, _default uint) uint {
	v := c.GetDefault(path, _default)
	return v.ToUint()
}

func (c *config) GetFloat32(path string, _default float32) float32 {
	v := c.GetDefault(path, _default)
	return v.ToFloat32()
}

func (c *config) GetFloat64(path string, _default float64) float64 {
	v := c.GetDefault(path, _default)
	return v.ToFloat64()
}

func (c *config) AllSettings() map[string]any {
	return c.content.content
}

// SetInt sets an int value at the specified path.
func (c *config) SetInt(path string, value int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}

// SetUint sets an uint value at the specified path.
func (c *config) SetUint(path string, value uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}

// SetFloat32 sets a float32 value at the specified path.
func (c *config) SetFloat32(path string, value float32) error {
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

// SetTime sets a time.Time value at the specified path.
func (c *config) SetTime(path string, value time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.content.Set(path, value)
}
