package benchmark

import (
	"github.com/yates-z/easel/logger"
	"github.com/yates-z/easel/logger/backend"
	"go.uber.org/zap"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

// BenchmarkDefaultLog-10
// 27848             41477 ns/op               5 B/op          0 allocs/op
// 28186             41704 ns/op               5 B/op          0 allocs/op
// 28327             41598 ns/op               5 B/op          0 allocs/op
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
		log.Infos("hello infos",
			logger.F("", strconv.Itoa(n)),
			logger.F("", strconv.Itoa(n)),
		)
	}
}

// BenchmarkSlog-10    	  319574	      4160 ns/op
// 84526             13487 ns/op              48 B/op          1 allocs/op
// 86983             13399 ns/op              48 B/op          1 allocs/op
// 84753             13311 ns/op              48 B/op          1 allocs/op
func BenchmarkSlog(b *testing.B) {
	s := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.Info("hello world", slog.Int("", n), slog.Int("a", n), slog.Int("b", n))

	}
}

// 70202             16450 ns/op             281 B/op          3 allocs/op
// 70435             16486 ns/op             280 B/op          3 allocs/op
// 70279             16521 ns/op             280 B/op          3 allocs/op
func BenchmarkZap(b *testing.B) {
	log, _ := zap.NewProduction()
	//defer log.Sync()

	sugar := log.Sugar()
	b.ResetTimer()
	n := 0
	for ; n < b.N; n++ {
		sugar.Info("hello world", n)
		//fmt.Println(n)
		//sugar.Infof("hello %d", n)
	}
	//fmt.Println(n)

}
