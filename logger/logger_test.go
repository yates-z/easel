package logger

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/logger/backend"
	"github.com/yates-z/easel/logger/buffer"
	"os"
	"runtime/debug"
	"testing"
)

func TestLog(t *testing.T) {
	l := NewLogger(
		WithLevel(WarnLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(InfoLevel|WarnLevel, " "),
		WithSeparator(ErrorLevel|FatalLevel|PanicLevel, "@@"),
		WithFields(AnyLevel,
			DatetimeField("2006-01-02 15:04:03").Key("datetime"),
			LevelField().Key("level").Upper().Prefix("[").Suffix("]"),
			CallerField(true, false).Key("file"),
			FuncNameField(true).Key("func"),
			MessageField().Key("msg"),
			Group("sys_info",
				CustomField(func(buf *buffer.Buffer) {
					buildInfo, _ := debug.ReadBuildInfo()
					buf.WriteString(buildInfo.GoVersion)
				}).Key("go_version"),
				Group("sys", CustomField(func(buf *buffer.Buffer) {
					buf.WriteInt(int64(os.Getpid()))
				}).Key("pid")),
			),
		),
		WithEncoders(AnyLevel, PlainEncoder, JSONEncoder, LogFmtEncoder),
	)
	l.Debug("hello debug")
	l.Info("hello info")
	l.Warn("hello warn")
	l.Error("hello error")

	//l.Fatal("hello fatal")
	//l.Panic("hello panic")
}

func TestLogf(t *testing.T) {
	l := NewLogger(
		WithLevel(InfoLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(InfoLevel|WarnLevel, " "),
		WithSeparator(ErrorLevel|FatalLevel|PanicLevel, "@@"),
		WithFields(AnyLevel,
			DatetimeField("2006-01-02 15:04:03").Key("datetime"),
			LevelField().Key("level").Upper().Prefix("[").Suffix("]"),
			CallerField(true, true).Key("caller"),
			MessageField().Key("msg"),
			Group("sys_info",
				CustomField(func(buf *buffer.Buffer) {
					buildInfo, _ := debug.ReadBuildInfo()
					buf.WriteString(buildInfo.GoVersion)
				}).Key("go_version"),
				Group("sys", CustomField(func(buf *buffer.Buffer) {
					buf.WriteInt(int64(os.Getpid()))
				}).Key("pid")),
			),
		),
		WithEncoders(AnyLevel, PlainEncoder, JSONEncoder, LogFmtEncoder),
	)

	l.Debugf("hello %s %d", "debugf", 1000)
	l.Infof("hello %s %d", "infof", 1001)
	l.Warnf("hello %s %.2f", "warnf", 1002)
	l.Errorf("hello %s %d", "errorf", 1003)
	//l.Fatalf("hello %s %d", "fatalf", 1000)
	//l.Panicf("hello %s %d", "panicf", 1000)
}

func TestLogWithColor(t *testing.T) {
	l := NewLogger(
		WithLevel(DebugLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(InfoLevel|WarnLevel, " "),
		WithSeparator(ErrorLevel|FatalLevel|PanicLevel, " "),
		WithFields(AnyLevel,
			DatetimeField("2006-01-02 15:04:03").Key("datetime").Color(Yellow),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Red),
		),
		WithFields(AnyLevel^ErrorLevel^FatalLevel^PanicLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]"),
		),
		WithFields(AnyLevel,
			CallerField(true, true).Key("caller").Color(Black),
			MessageField().Key("msg").Background(Blue),
			Group("sys_info",
				CustomField(func(buf *buffer.Buffer) {
					buildInfo, _ := debug.ReadBuildInfo()
					buf.WriteString(buildInfo.GoVersion)
				}).Key("go_version"),
				Group("sys", CustomField(func(buf *buffer.Buffer) {
					buf.WriteInt(int64(os.Getpid()))
				}).Key("pid")),
			),
		),
		WithEncoders(AnyLevel, PlainEncoder, JSONEncoder, LogFmtEncoder),
	)
	l.Debug("hello debug")
	l.Info("hello info")
	l.Warn("hello warn")
	l.Error("hello error")

	//l.Fatal("hello fatal")
	//l.Panic("hello panic")
}

func TestDefaultLog(t *testing.T) {
	Debug("hello", "debug1", "debug2", 6666)
	Info("hello info")
	Warn("hello warn")
	Error("hello error")

	Debugf("hello %s %d", "debugf", 1000)
	Infof("hello %s %d", "infof", 1001)
	Warnf("hello %s %.2f", "warnf", 1002)
	Errorf("hello %s %d", "errorf", 1003)

	Context(context.Background()).Debug("hello debug")

}

func TestLogs(t *testing.T) {
	l := NewLogger(
		WithLevel(DebugLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(AnyLevel, "    "),
		WithFields(AnyLevel,
			DatetimeField("2006/01/02 15:04:03").Key("datetime"),
		),
		WithFields(DebugLevel|InfoLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Green),
		),
		WithFields(WarnLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Yellow),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField().Key("level").Upper().Prefix("[").Suffix("]").Color(Red),
		),
		WithFields(AnyLevel,
			MessageField().Key("msg"),
			CallerField(true, true).Key("caller"),
		),
		WithEncoders(AnyLevel, PlainEncoder, JSONEncoder, LogFmtEncoder),
	)
	l.Debugs("hello debug", F("test", "this is a test"))
	l.Debugs("hello debug",
		Group("TEST",
			F("test", "6").Color(Yellow),
			F("test", "7").Color(Blue),
			F("test", "7").Color(Blue),
		),
		F("tset", "8").Color(Green),
	)
}

func TestGetBackends(t *testing.T) {
	l := NewLogger(
		WithBackends(DebugLevel, backend.OSBackend().Build()),
		WithBackends(InfoLevel, backend.DefaultFileBackend().Build()),
		WithBackends(WarnLevel, backend.DefaultFileBackend().Build()),
	)
	fmt.Println(l.Backends())
	l = NewLogger(
		WithLevel(DebugLevel),
		WithBackends(DebugLevel, backend.OSBackend().Build()),
		WithBackends(InfoLevel, backend.DefaultFileBackend().Build()),
		WithBackends(WarnLevel, backend.DefaultFileBackend().Build()),
	)
	fmt.Println(l.Backends())
}
