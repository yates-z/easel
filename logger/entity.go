package logger

import (
	"errors"
	"fmt"
	"github.com/yates-z/easel/logger/backend"
)

type logEntity struct {
	level LogLevel
	opts  *entityOptions
}

func (e *logEntity) preLog(msg *string, withColor bool) []byte {
	var logs []byte
	for _, encoder := range e.opts.encoders {
		log := encoder.Encode(newRecord(e, msg, withColor))
		for count := e.opts.skipLines; count >= 0; count-- {
			log += "\n"
		}
		logs = append(logs, []byte(log)...)
	}
	return logs
}

func (e *logEntity) log(msg *string) (errs []error, available map[backend.Backend]struct{}) {
	var coloredLogs []byte
	logs := e.preLog(msg, false)
	available = make(map[backend.Backend]struct{})
	for b := range e.opts.backends {
		var err error
		if b.AllowANSI() {
			if len(coloredLogs) == 0 {
				coloredLogs = e.preLog(msg, true)
			}
			_, err = b.Write(coloredLogs)
		} else {
			_, err = b.Write(logs)
		}
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("%T writer error: %s.", b, err)))
		} else {
			available[b] = struct{}{}
		}
	}
	return
}

func (e *logEntity) handleError(errs []error, available map[backend.Backend]struct{}) {
	if len(errs) == 0 || !e.level.Eq(ErrorLevel) {
		return
	}
	// broadcast errors
	var msg string
	for _, err := range errs {
		msg += err.Error()
	}
	entity := e.copy()
	entity.opts.backends = available
	entity.log(&msg)

}

func (e *logEntity) copy() *logEntity {
	entity := &logEntity{
		level: e.level,
		opts: &entityOptions{
			separator: e.opts.separator,
			skipLines: e.opts.skipLines,
			fields:    e.opts.fields,
			encoders:  e.opts.encoders,
			backends:  e.opts.backends,
		},
	}
	return entity
}
