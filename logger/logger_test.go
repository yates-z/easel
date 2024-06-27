package logger

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/logger/backend"
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
			DatetimeField("2006-01-02 15:04:03").Key("datetime").Build(),
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Build(),
			CallerField(true, false).Key("file").Build(),
			FuncNameField(true).Key("func").Build(),
			MessageField().Key("msg").Build(),
			Group("sys_info",
				CustomField(func() string {
					buildInfo, _ := debug.ReadBuildInfo()
					return buildInfo.GoVersion
				}).Key("go_version").Build(),
				Group("sys", CustomField(func() string {
					return fmt.Sprintf("%d", os.Getpid())
				}).Key("pid").Build()).Build(),
			).Build(),
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
			DatetimeField("2006-01-02 15:04:03").Key("datetime").Build(),
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Build(),
			CallerField(true, true).Key("caller").Build(),
			MessageField().Key("msg").Build(),
			Group("sys_info",
				CustomField(func() string {
					buildInfo, _ := debug.ReadBuildInfo()
					return buildInfo.GoVersion
				}).Key("go_version").Build(),
				Group("sys", CustomField(func() string {
					return fmt.Sprintf("%d", os.Getpid())
				}).Key("pid").Build()).Build(),
			).Build(),
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
			DatetimeField("2006-01-02 15:04:03").Key("datetime").Color(Yellow).Build(),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Red).Build(),
		),
		WithFields(AnyLevel^ErrorLevel^FatalLevel^PanicLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Build(),
		),
		WithFields(AnyLevel,
			CallerField(true, true).Key("caller").Color(Black).Background(Blue).Build(),
			MessageField().Key("msg").Build(),
			Group("sys_info",
				CustomField(func() string {
					buildInfo, _ := debug.ReadBuildInfo()
					return buildInfo.GoVersion
				}).Key("go_version").Build(),
				Group("sys", CustomField(func() string {
					return fmt.Sprintf("%d", os.Getpid())
				}).Key("pid").Build()).Build(),
			).Build(),
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
	Debug("hello debug")
	Info("hello info")
	Warn("hello warn")
	Error("hello error")

	Debugf("hello %s %d", "debugf", 1000)
	Infof("hello %s %d", "infof", 1001)
	Warnf("hello %s %.2f", "warnf", 1002)
	Errorf("hello %s %d", "errorf", 1003)

	Context(context.Background()).Debug("hello debug")

}
