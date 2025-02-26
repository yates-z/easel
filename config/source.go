package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Source interface {
	Load() (*Content, error)
}

var _ Source = (*file)(nil)

// file is a source that loads configuration from a file.
type file struct {
	path string
}

func NewFile(path string) *file {
	return &file{path: path}
}

func (f *file) Load() (*Content, error) {
	p := strings.Split(f.path, ".")
	if len(p) <= 1 {
		return nil, errors.New("invalid file path")
	}

	bytes, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	ext := p[len(p)-1]
	switch ext {
	case "yaml", "yml":
		return parseYAML(bytes)
	case "json":
		return parseJSON(bytes)
	default:
		return nil, errors.New("unsupported file extension")
	}
}

func parseYAML(bytes []byte) (*Content, error) {
	c := NewContent()
	err := yaml.Unmarshal(bytes, &c.content)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	return c, nil
}

func parseJSON(bytes []byte) (*Content, error) {
	c := NewContent()
	err := json.Unmarshal(bytes, &c.content)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	return c, nil
}

var _ Source = (*environment)(nil)

// environment is a source that loads environment variables.
type environment struct {
	prefixes []string
}

func NewEnviron(prefixes ...string) *environment {
	return &environment{prefixes: prefixes}
}

// Load implements Source.
func (e *environment) Load() (*Content, error) {
	c := NewContent()

	// handle environment variables.
	for _, env := range os.Environ() {

		// split the key and value.
		kv := strings.SplitN(env, "=", 2)

		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		if len(e.prefixes) > 0 {
			p, ok := matchPrefix(e.prefixes, key)
			if !ok || len(p) == len(key) {
				continue
			}
			// trim prefix
			key = strings.TrimPrefix(key, p)
			key = strings.TrimPrefix(key, "_")
		}

		if len(key) != 0 {
			err := c.Set(key, value)
			if err != nil {
				return nil, fmt.Errorf("handle environ variables %s failed: %w", key, err)
			}
		}
	}
	return c, nil
}

func matchPrefix(prefixes []string, s string) (string, bool) {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}
