package test

import (
	"testing"

	"github.com/yates-z/easel/config"
)

func TestConfig(t *testing.T) {
	config.Load(config.DefaultConfigPath)
	host := config.Get("databases.main.port")
	println(host.ToInt())

	version := config.Get("databases.main.version[0][1]")
	println(version.ToString())
}

func TestSetConfig(t *testing.T) {
	config.Load(config.DefaultConfigPath)
	// test setting a new key
	config.SetInt("mytest.main.port", 3306)
	port := config.Get("mytest.main.port")
	println(port.ToInt() == 3306)

	// test setting a new key with an array index
	err := config.SetString("mytest.version[0]", "1.0.0")
	if err != nil {
		println(err.Error())
		return
	}
	version := config.Get("mytest.version[0]")
	println(version.ToString() == "1.0.0")
}

func TestSetConfig2(t *testing.T) {
	config.Load(config.DefaultConfigPath)
	// test setting an non-existing key
	err := config.SetString("mytest.version[1][0]", "2.0.0")
	if err != nil {
		println(err.Error())
		return
	}
	version := config.Get("mytest.version[1][0]")
	println(version.ToString() == "2.0.0")
	env_gopath := config.Get("GOPATH")
	println(env_gopath.ToString())
}

func TestJsonConfig(t *testing.T) {
	config.Load("config.json")
	namespace := config.Get("namespace")
	println(namespace.ToInt() == 1111)
	filePath := config.Get("filePath")
	println(filePath.ToString())
	phone := config.Get("address[0].phone")
	println(phone.ToString())

	band := config.Get("myCar.band")
	println(band.ToString())

}
