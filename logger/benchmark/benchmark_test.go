package benchmark

import (
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/logger/backend"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

// BenchmarkDefaultLog-10
// 77700             15293 ns/op               5 B/op          0 allocs/op
// 74779             15248 ns/op               5 B/op          0 allocs/op
// 75206             15215 ns/op               5 B/op          0 allocs/op
func BenchmarkDefault(b *testing.B) {
	log := logger.NewLogger(
		logger.WithLevel(logger.InfoLevel),
		logger.WithBackends(logger.AnyLevel, backend.OSBackend().Build()),
		logger.WithSeparator(logger.AnyLevel, " "),
		logger.WithFields(logger.AnyLevel,
			logger.DatetimeField("2006/01/02 15:04:05"),
			logger.LevelField().Upper(),
			logger.MessageField(),
			//logger.TimeField(logger.UnixMilli).Background(logger.Green),
			//logger.FuncNameField(true),
			//logger.CallerField(true, true),
		),
		logger.WithEncoders(logger.AnyLevel, logger.LogFmtEncoder),
	)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		//log.Info("Hello World", n)
		log.Infos("hello world",
			logger.Group("TEST",
				logger.F("test", "6"),
				logger.F("test", "7"),
				logger.F("test", "7"),
			),
			logger.F("test", strconv.Itoa(n)),
		)
	}
}

// BenchmarkSlog-10    	  319574	      4160 ns/op
// 74466             15896 ns/op             533 B/op          8 allocs/op
// 71007             15828 ns/op             533 B/op          8 allocs/op
// 74871             15882 ns/op             533 B/op          8 allocs/op
func BenchmarkSlog(b *testing.B) {
	s := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.Info("hello world",
			slog.Group("TEST",
				slog.String("test", "6"),
				slog.String("test", "7"),
				slog.String("test", "7"),
			),
			slog.String("test", strconv.Itoa(n)),
		)
	}
}
