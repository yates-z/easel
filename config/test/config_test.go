package test

import (
	"testing"

	"github.com/yates-z/easel/config"
)

func TestConfig(t *testing.T) {
	config.LoadConfig("config.yml")
	host := config.Get("databases.main.port")
	println(host.ToInt())

	version := config.Get("databases.main.version[1]")
	println(version.ToString())
}
