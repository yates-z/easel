package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"testing"
)

type Test struct {
	Test1 string
	Test2 string
	Test3 string
}

func (a Test) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("test", a.Test1)
	enc.AddString("test", a.Test2)
	enc.AddString("test", a.Test3)
	return nil
}

// 29853             39203 ns/op             245 B/op          5 allocs/op
// 29844             39290 ns/op             245 B/op          5 allocs/op
// 29671             39341 ns/op             245 B/op          5 allocs/op
func BenchmarkZap(b *testing.B) {
	//defer log.Sync()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "datetime",
		LevelKey:       "level",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 创建core
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	log := zap.New(core)
	test := Test{
		Test1: "6",
		Test2: "7",
		Test3: "7",
	}

	b.ResetTimer()
	n := 0
	for ; n < b.N; n++ {
		log.Info("hello world",
			zap.Object("TEST", test),
			zap.String("test", strconv.Itoa(n)),
		)
		//fmt.Println(n)
		//sugar.Infof("hello %d", n)
	}
	//fmt.Println(n)

}
