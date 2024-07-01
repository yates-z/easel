package logger

import (
	"github.com/yates-z/easel/logger/buffer"
)

var (
	PlainEncoder  = &plainEncoder{}
	JSONEncoder   = &jsonEncoder{}
	LogFmtEncoder = &logFmtEncoder{}
)

type Record struct {
	entity    *logEntity
	msg       string
	withColor bool
	buf       *buffer.Buffer
}

func newRecord(entity *logEntity, msg string, b *buffer.Buffer) Record {
	return Record{
		entity:    entity,
		msg:       msg,
		withColor: false,
		buf:       b,
	}
}

func (r *Record) free() {
	r.msg = ""
	r.buf.Free()
}

func newRecordWithColor(entity *logEntity, msg string, b *buffer.Buffer) Record {
	return Record{
		entity:    entity,
		msg:       msg,
		withColor: true,
		buf:       b,
	}
}

type Encoder interface {
	Encode(r Record) *buffer.Buffer
}

type plainEncoder struct {
}

func (e *plainEncoder) encode(field LogField, r Record) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				_, _ = r.buf.WriteString(r.entity.opts.separator)
			}
			e.encode(child, r)
		}
		return
	}

	if field.Kind() == LevelFieldType {
		field.decorate(r.buf, r.entity.level.String(), r.withColor)
	} else if field.Kind() == MessageFieldType {
		field.decorate(r.buf, r.msg, r.withColor)
	} else {
		field.decorate(r.buf, "", r.withColor)
	}
}

func (e *plainEncoder) Encode(r Record) *buffer.Buffer {
	for index, field := range r.entity.opts.fields {
		if index != 0 {
			_, _ = r.buf.WriteString(r.entity.opts.separator)
		}

		e.encode(field, r)
	}
	return r.buf
}

type jsonEncoder struct{}

func (e *jsonEncoder) encode(field LogField, r Record) {

	if len(field.Children()) > 0 {
		key := field.Key()
		_, _ = r.buf.WriteString(`"` + key + `": {`)
		for index, child := range field.Children() {
			if index != 0 {
				_, _ = r.buf.WriteString(", ")
			}
			e.encode(child, r)
		}
		_ = r.buf.WriteByte('}')
		return
	}

	_, _ = r.buf.WriteString(`"` + field.Key() + `":"`)

	if field.Kind() == LevelFieldType {
		field.decorate(r.buf, r.entity.level.String(), r.withColor)
	} else if field.Kind() == MessageFieldType {
		field.decorate(r.buf, r.msg, r.withColor)
	} else {
		field.decorate(r.buf, "", r.withColor)
	}
	_ = r.buf.WriteByte('"')
	//value = field.Decorate(value, r.color)

}

func (e *jsonEncoder) Encode(r Record) *buffer.Buffer {
	_ = r.buf.WriteByte('{')
	for index, field := range r.entity.opts.fields {
		if index != 0 {
			_, _ = r.buf.WriteString(", ")
		}
		e.encode(field, r)
	}
	_ = r.buf.WriteByte('}')

	return r.buf
}

type logFmtEncoder struct{}

func (e *logFmtEncoder) encode(field LogField, key []byte, r Record) {

	if field.Kind() == GroupFieldType {
		for index, child := range field.Children() {
			if index != 0 {
				_ = r.buf.WriteByte(' ')
			}
			subKey := key
			subKey = append(subKey, '.')
			subKey = append(subKey, child.Key()...)
			e.encode(child, subKey, r)
		}
		return
	}

	_, _ = r.buf.Write(key)
	_ = r.buf.WriteByte('=')

	if field.Kind() == LevelFieldType {
		field.decorate(r.buf, r.entity.level.String(), r.withColor)
	} else if field.Kind() == MessageFieldType {
		_ = r.buf.WriteByte('"')
		field.decorate(r.buf, r.msg, r.withColor)
		_ = r.buf.WriteByte('"')
	} else {
		field.decorate(r.buf, "", r.withColor)
	}

}

func (e *logFmtEncoder) Encode(r Record) *buffer.Buffer {
	for index, field := range r.entity.opts.fields {
		if index != 0 {
			_ = r.buf.WriteByte(' ')
		}
		key := []byte(field.Key())
		e.encode(field, key, r)
	}
	return r.buf
}
