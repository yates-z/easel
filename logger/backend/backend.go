package backend

import "io"

type Backend interface {
	io.Writer
	//Sync flushes buffered logs
	Sync() error
	//AllowANSI determines if allow to show colorful log
	AllowANSI() bool
}
