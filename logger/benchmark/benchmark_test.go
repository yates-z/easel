package benchmark

import (
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/logger/backend"
	"go.uber.org/zap"
	"log/slog"
	"os"
	"testing"
)

// BenchmarkDefaultLog-10
// 34527             34066 ns/op               0 B/op          0 allocs/op
// 33933             34096 ns/op               0 B/op          0 allocs/op
// 34754             34128 ns/op               0 B/op          0 allocs/op
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
			logger.CallerField(true, true),
		),
		logger.WithEncoders(logger.AnyLevel, logger.LogFmtEncoder),
	)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Info("Hello World")
		//log.Warns("hello warns",
		//	logger.Group("TEST",
		//		logger.F("test", "6").Color(logger.Green),
		//		logger.F("test", "7").Color(logger.Yellow),
		//		logger.F("test", "7").Color(logger.Yellow),
		//	),
		//	logger.F("tset", "8"),
		//)
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

// 38714661               157.7 ns/op             2 B/op          0 allocs/op
// 38072233               150.7 ns/op             2 B/op          0 allocs/op
// 34610774               152.1 ns/op             2 B/op          0 allocs/op
func BenchmarkZap(b *testing.B) {
	log, _ := zap.NewProduction()
	defer log.Sync()

	sugar := log.Sugar()
	sugar = sugar.WithOptions(zap.AddCaller())
	b.ResetTimer()
	n := 0
	for ; n < b.N; n++ {
		sugar.Info("hello world")
		//sugar.Infof("hello %d", n)
	}
	//fmt.Println(n)

}
