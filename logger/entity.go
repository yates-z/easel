package logger

import (
	"errors"
	"fmt"
	"github.com/yates-z/easel/logger/backend"
	"github.com/yates-z/easel/logger/buffer"
)

type logEntity struct {
	level    LogLevel
	opts     *entityOptions
	backends []backend.Backend
}

func (e *logEntity) preLog(record Record) {

	for _, encoder := range e.opts.encoders {
		buf := encoder.Encode(record)
		for count := e.opts.skipLines; count >= 0; count-- {
			_ = buf.WriteByte('\n')
		}
	}
	return
}

func (e *logEntity) log(msg string) (errs []error) {
	buf := buffer.New()
	coloredBuf := buffer.New()
	defer buf.Free()
	defer coloredBuf.Free()

	for _, b := range e.backends {
		var err error
		if b.AllowANSI() {
			if coloredBuf.Len() == 0 {
				record := newRecordWithColor(e, msg, coloredBuf)
				e.preLog(record)
			}
			_, err = b.Write(*coloredBuf)
		} else {
			if buf.Len() == 0 {
				record := newRecord(e, msg, buf)
				e.preLog(record)
			}
			_, err = b.Write(*buf)
		}
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("%T writer error: %s.", b, err)))
		}
	}
	return
}

func (e *logEntity) handleError(errs []error) {
	if len(errs) == 0 || !e.level.Eq(ErrorLevel) {
		return
	}
	// broadcast errors
	var msg string
	for _, err := range errs {
		msg += err.Error()
	}
	//entity := e.copy()
	//e.log(&msg)

}

func (e *logEntity) copy(target *logEntity) {
	e.level = target.level
	e.opts.separator = target.opts.separator
	e.opts.skipLines = target.opts.skipLines
	e.opts.fields = append(e.opts.fields[:0], target.opts.fields...)
	e.opts.encoders = append(e.opts.encoders[:0], target.opts.encoders...)

	for name, b := range target.opts.backends {
		e.opts.backends[name] = b
		e.backends = append(e.backends, b)
	}
}

func (e *logEntity) clear() {
	e.level = 0
	e.opts.separator = ""
	e.opts.skipLines = 0
	e.opts.fields = e.opts.fields[:0]
	e.opts.encoders = e.opts.encoders[:0]
	for key := range e.opts.backends {
		delete(e.opts.backends, key)
	}
	e.backends = e.backends[:0]
}
