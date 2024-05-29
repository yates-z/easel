package logger

type LogLevel uint8

const (
	DebugLevel = iota
	// InfoLevel is the default log level
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}
	return ""
}

func (l LogLevel) Enabled(level LogLevel) bool {
	return level >= l
}

func (l LogLevel) Eq(level LogLevel) bool {
	return level == l
}

func (l LogLevel) Enum() []LogLevel {
	var levels []LogLevel
	var level = l
	for ; level <= PanicLevel; level++ {
		levels = append(levels, level)
	}
	return levels
}
