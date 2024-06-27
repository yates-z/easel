package logger

type LogLevel uint8

const (
	DebugLevel LogLevel = 1
	// InfoLevel is the default log level
	InfoLevel    LogLevel = 1 << 1
	WarnLevel    LogLevel = 1 << 2
	ErrorLevel   LogLevel = 1 << 3
	FatalLevel   LogLevel = 1 << 4
	PanicLevel   LogLevel = 1 << 5
	HighestLevel          = PanicLevel

	AnyLevel = DebugLevel | InfoLevel | WarnLevel | ErrorLevel | FatalLevel | PanicLevel
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
	return level >= l && level <= HighestLevel && (level%2 == 0 || level == DebugLevel)
}

func (l LogLevel) Eq(level LogLevel) bool {
	return level == l
}

// EnumIncremental returns LogLevels which have higher level.
func (l LogLevel) EnumIncremental() []LogLevel {
	var levels []LogLevel
	var level = l
	for ; level <= HighestLevel; level = level << 1 {
		levels = append(levels, level)
	}
	return levels
}

func (l LogLevel) Contains(level LogLevel) bool {
	result := uint8(l) & uint8(level)
	return result == uint8(level)
}

func (l LogLevel) Enum() []LogLevel {
	var levels []LogLevel
	var level LogLevel
	for level = DebugLevel; level <= HighestLevel; level = level << 1 {
		if l.Contains(level) {
			levels = append(levels, level)
		}
	}
	return levels
}
