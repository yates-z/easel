package config

import (
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

// LoadConfig reads the YAML configuration file and loads it into the values map
func LoadConfig(filePath string) {
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

// Set sets a configuration value by its key, supporting nested keys using dot notation (e.g., "database.host").
// Note: The Set method does not allow setting values for non-existent keys, If necessary, please use the Add method.
// TODO: Implement the Set method
func Set(path string, value any) error {
	config.mu.Lock()
	defer config.mu.Unlock()

	return nil
}

// Add adds a configuration value by its key, supporting nested keys using dot notation (e.g., "database.host").
// Note: The Add method does not allow setting values for existing keys. If necessary, please use the Set method.
// TODO: Implement the Add method
func Add(path string, value any) error {
	return nil
}

// overrideWithEnvVariables overrides configuration values with environment variables if they exist
func overrideWithEnvVariables() {

	// handle environment variables.
	for _, env := range os.Environ() {

		// split the key and value.
		kv := strings.SplitN(env, "=", 2)
		println(kv)
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		//
		if err := Set(key, value); err != nil {
			logger.Errorf("handle environ variables %s failed: %w", key, err)
		}
	}
}
