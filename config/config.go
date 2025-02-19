package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/yates-z/easel/core/variant"
	"github.com/yates-z/easel/logger"
)

const DefaultConfigPath = "config.yaml"

type Config struct {
	mu   sync.RWMutex
	data map[string]any
}

var config *Config = &Config{}

// LoadFile reads the YAML configuration file and loads it into the values map
func LoadFile(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config.data)
	if err != nil {
		logger.Fatalf("Error parsing config file: %v", err)
	}

	overrideWithEnvVariables()
}

// Get retrieves a configuration value by its key, supporting nested keys using dot notation (e.g., "database.host")
func Get(path string) variant.Variant {
	config.mu.RLock()
	defer config.mu.RUnlock()

	parts := parsePath(path)

	var current any = config.data

	for _, part := range parts {
		if part.kind == Index {
			index, error := strconv.Atoi(part.value)
			if error != nil {
				return variant.Nil
			}
			currentSlice, ok := current.([]any)
			if !ok || index >= len(currentSlice) || index < 0 {
				return variant.Nil
			}
			current = currentSlice[index]
		} else if part.kind == Key {
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return variant.Nil
			}
			val, exists := currentMap[part.value]
			if !exists {
				return variant.Nil
			}
			current = val
		}
	}

	return variant.New(current)
}

// set sets a configuration value by its key, supporting nested keys using dot notation (e.g., "database.host").
func set(path string, value any) error {
	config.mu.Lock()
	defer config.mu.Unlock()

	parts := parsePath(path)

	var parent any
	var current any = config.data
	for i, part := range parts {
		// fmt.Printf("parent: %+v, current: %+v\n\n", parent, current)
		if part.kind == Index {
			index, error := strconv.Atoi(part.value)
			if error != nil {
				return error
			}

			currentSlice, ok := current.([]any)
			if !ok {
				return errors.New("invalid type")
			}

			// expand the length of the slice.
			if index > len(currentSlice) || index < 0 {
				return errors.New("index out of range")
			}
			if index == len(currentSlice) {
				newSlice := make([]any, index+1)
				copy(newSlice, currentSlice)
				currentSlice = newSlice

				if part.parent.kind == Key {
					parentMap, _ := parent.(map[string]any)
					parentMap[part.parent.value] = newSlice
				} else {
					parentSlice, _ := parent.([]any)
					parentIndex, _ := strconv.Atoi(part.parent.value)
					parentSlice[parentIndex] = newSlice
				}
			}

			// set the value if it is a leaf node.
			if part.isLeaf {
				currentSlice[index] = value
				return nil
			}

			if currentSlice[index] == nil {
				if parts[i+1].kind == Index {
					currentSlice[index] = make([]any, 0)
				} else {
					currentSlice[index] = make(map[string]any)
				}
			}
			parent = currentSlice
			current = currentSlice[index]
		} else if part.kind == Key {
			currentMap, ok := current.(map[string]any)
			if !ok {
				return errors.New("invalid type")
			}

			// set the value if it is a leaf node.
			if part.isLeaf {
				currentMap[part.value] = value
				return nil
			}

			// create a new map/slice if the key does not exist.a.b.c
			if _, exists := currentMap[part.value]; !exists {
				if parts[i+1].kind == Index {
					currentMap[part.value] = make([]any, 0)
				} else {
					currentMap[part.value] = make(map[string]any)
				}
			}
			parent = currentMap
			current = currentMap[part.value]
		}
	}

	return nil
}

func SetInt(path string, value int) error {
	return set(path, value)
}

func SetFloat64(path string, value float64) error {
	return set(path, value)
}

func SetString(path string, value string) error {
	return set(path, value)
}

func SetBool(path string, value bool) error {
	return set(path, value)
}

// overrideWithEnvVariables overrides configuration values with environment variables if they exist
func overrideWithEnvVariables() {

	// handle environment variables.
	for _, env := range os.Environ() {

		// split the key and value.
		kv := strings.SplitN(env, "=", 2)

		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		//
		if err := set(key, value); err != nil {
			logger.Errorf("handle environ variables %s failed: %w", key, err)
		}
	}
}
