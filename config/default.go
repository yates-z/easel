package config

import "github.com/yates-z/easel/core/variant"

const DefaultConfigPath = "config.yaml"

var defaultConfig *config = &config{}

func Get(path string) (variant.Variant, bool) {
	return defaultConfig.Get(path)
}

func SetInt(path string, value int) error {
	return defaultConfig.SetInt(path, value)
}

func SetFloat64(path string, value float64) error {
	return defaultConfig.SetFloat64(path, value)
}

func SetString(path string, value string) error {
	return defaultConfig.SetString(path, value)
}

func SetBool(path string, value bool) error {
	return defaultConfig.SetBool(path, value)
}
