package logger

type LogMode uint8

const (
	DebugMode = iota
	TestMode
	PreProductMode
	ProductMode
)
