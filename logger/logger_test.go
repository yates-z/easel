package logger

import (
	"context"
	"fmt"
	"github.com/yates-z/easel/logger/backend"
	"log/slog"
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
			DatetimeField("2006/01/02 15:04:03").Key("datetime").Build(),
		),
		WithFields(DebugLevel|InfoLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Green).Build(),
		),
		WithFields(WarnLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Yellow).Build(),
		),
		WithFields(ErrorLevel|FatalLevel|PanicLevel,
			LevelField(true).Key("level").Upper().Prefix("[").Suffix("]").Color(Red).Build(),
		),
		WithFields(AnyLevel,
			MessageField().Key("msg").Build(),
			CallerField(true, true).Key("caller").Build(),
		),
		WithEncoders(AnyLevel, PlainEncoder, JSONEncoder, LogFmtEncoder),
	)
	l.Debugs("hello debug", F("test", "this is a test").Build())
	l.Debugs("hello debug",
		F("test", 6).Color(Yellow).Build(),
		F("test", 7).Color(Blue).Build(),
		F("tset", 8).Color(Green).Build(),
	)
}

// BenchmarkDefaultLog-10
// 450146             11364 ns/op            1547 B/op         53 allocs/op
// 478185             11418 ns/op            1547 B/op         53 allocs/op
// 467475             11527 ns/op            1547 B/op         53 allocs/op
func BenchmarkDefaultLog(b *testing.B) {
	l := NewLogger(
		WithLevel(InfoLevel),
		WithBackends(AnyLevel, backend.OSBackend().Build()),
		WithSeparator(AnyLevel, " "),
		WithFields(AnyLevel,
			DatetimeField("2006-01-02 15:04:03").Key("datetime").Build(),
			LevelField(true).Key("level").Upper().Build(),
			MessageField().Key("msg").Build(),
		),
		WithEncoders(AnyLevel, PlainEncoder),
	)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		l.Info("hello world")
	}
}

// BenchmarkSlog-10    	  319574	      4160 ns/op
// 477282             11529 ns/op               0 B/op          0 allocs/op
// 476253             11671 ns/op               0 B/op          0 allocs/op
// 455965             11652 ns/op               0 B/op          0 allocs/op
func BenchmarkSlog(b *testing.B) {
	s := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.Info("hello world")
	}
}
