package logger

import (
	"fmt"
	"github.com/yates-z/easel/logger/backend"
	"os"
	"runtime/debug"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(
		WithLevel(InfoLevel),
		WithSkipLines(1),
		WithBackends(backend.OSBackend().Build(), backend.DefaultFileBackend().Build()),
		WithFields(
			LevelField("level").Upper(true).Build(),
			DatetimeField("datetime").Build(),
			ShortFileField("file").Build(),
			MessageField("msg").Build(),
		),
		WithSeparator(" "),
		WithEncoders(PlainEncoder(), JsonEncoder()),
	)
	logger.Log(DebugLevel, "this is a test")
	logger.Log(WarnLevel, "this is a test", 123, 456)
	logger.Logf(ErrorLevel, "this is a %s test", "error")
	logger.Logf(PanicLevel, "fatal..... %d", 666)
}

func TestColor(t *testing.T) {
	logger := NewLogger(
		WithLevel(InfoLevel),
		WithSkipLines(1),
		WithBackends(
			backend.OSBackend().Build(),
			backend.DefaultFileBackend().Filename("log.log").Build(),
		),
		WithFields(
			LevelField("level").Upper(true).Color(Red).Background(Blue).Build(),
			DatetimeField("datetime").Color(Green).Build(),
			ShortFileField("file").Color(Black).Background(Magenta).Build(),
			MessageField("msg").Background(Yellow).Build(),
			Group("sys_info",
				CustomField("go_version").Handle(func() string {
					buildinfo, _ := debug.ReadBuildInfo()
					return buildinfo.GoVersion
				}).Color(Red).Build(),
				Group("sys", CustomField("pid").Handle(func() string {
					return fmt.Sprintf("%d", os.Getpid())
				}).Build()).Build(),
			).Build(),
		),
		WithSeparator("   "),
		WithEncoders(JsonEncoder(), PlainEncoder(), LogFmtEncoder()),
	)
	logger.Log(DebugLevel, "this is a test")
	logger.Log(WarnLevel, "this is a test", 123, 456, nil)
	logger.Logf(ErrorLevel, "this is a %s test", "error")
}

func TestGroup(t *testing.T) {
	logger := NewLogger(
		WithLevel(InfoLevel),
		WithBackends(backend.OSBackend().Build(), backend.DefaultFileBackend().Build()),
		WithFields(
			LevelField("level").Upper(true).Build(),
			DatetimeField("datetime").Build(),
			ShortFileField("file").Build(),
			MessageField("msg").Build(),
			Group("sys_info",
				CustomField("go_version").Handle(func() string {
					buildinfo, _ := debug.ReadBuildInfo()
					return buildinfo.GoVersion
				}).Build(),
				Group("sys", CustomField("pid").Handle(func() string {
					return fmt.Sprintf("%d", os.Getpid())
				}).Build()).Build(),
			).Build(),
		),
		WithSeparator(" "),
		WithEncoders(PlainEncoder(), JsonEncoder(), LogFmtEncoder()),
	)
	logger.Log(DebugLevel, "this is a test")
	logger.Log(WarnLevel, "this is a test", 123, 456)
	logger.Logf(ErrorLevel, "this is a %s test", "error")
}
