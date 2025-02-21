package test

import (
	"testing"

	"github.com/yates-z/easel/config"
)

func TestConfig(t *testing.T) {
	c := config.New(config.WithSource(config.NewFile("config.yml")))
	host, _ := c.Get("databases.main.port")
	println(host.ToInt())

	version, _ := c.Get("databases.main.version[0][1]")
	println(version.ToString())
}

func TestSetConfig(t *testing.T) {
	config.New(config.WithSource(config.NewFile("config.yml")))

	// test setting a new key
	config.SetInt("mytest.main.port", 3306)
	port, _ := config.Get("mytest.main.port")
	println(port.ToInt() == 3306)

	// test setting a new key with an array index
	err := config.SetString("mytest.version[0]", "1.0.0")
	if err != nil {
		println(err.Error())
		return
	}
	version, _ := config.Get("mytest.version[0]")
	println(version.ToString())
}

func TestSetConfig2(t *testing.T) {
	config.New(config.WithSource(config.NewFile("config.yml")))
	// test setting an non-existing key
	err := config.SetString("mytest.version[1][0]", "2.0.0")
	if err != nil {
		println(err.Error())
		return
	}
	version, _ := config.Get("mytest.version[1][0]")
	println(version.ToString())
}
