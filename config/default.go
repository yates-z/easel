package config

import "github.com/yates-z/easel/core/variant"

const DefaultConfigPath = "config.yml"

var defaultConfig *config

func Load(path string) {
	defaultConfig = New(
		WithSource(NewFile(path)),
		WithSource(NewEnviron()),
	)
}

func Get(path string) variant.Variant {
	return defaultConfig.Get(path)
}

func GetDefault(path string, _default any) variant.Variant {
	return defaultConfig.GetDefault(path, _default)
}

func GetBool(path string, _default bool) bool {
	return defaultConfig.GetBool(path, _default)
}

func GetString(path string, _default string) string {
	return defaultConfig.GetString(path, _default)
}

func GetInt(path string, _default int) int {
	return defaultConfig.GetInt(path, _default)
}

func GetUint(path string, _default uint) uint {
	return defaultConfig.GetUint(path, _default)
}

func GetFloat32(path string, _default float32) float32 {
	return defaultConfig.GetFloat32(path, _default)
}

func GetFloat64(path string, _default float64) float64 {
	return defaultConfig.GetFloat64(path, _default)
}

func SetInt(path string, value int) error {
	return defaultConfig.SetInt(path, value)
}

func SetUint(path string, value uint) error {
	return defaultConfig.SetUint(path, value)
}

func SetFloat32(path string, value float32) error {
	return defaultConfig.SetFloat32(path, value)
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
