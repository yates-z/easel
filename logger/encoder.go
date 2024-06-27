package logger

import "fmt"

var (
	PlainEncoder  = &plainEncoder{}
	JSONEncoder   = &jsonEncoder{}
	LogFmtEncoder = &logFmtEncoder{}
)

type record struct {
	entity *logEntity
	msg    *string
	output string
	color  bool
}

func newRecord(entity *logEntity, msg *string, color bool) *record {
	return &record{
		entity: entity,
		msg:    msg,
		output: "",
		color:  color,
	}
}

type Encoder interface {
	Encode(r *record) string
}

type plainEncoder struct {
}

func (e *plainEncoder) encode(field LogField, r *record) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				r.output += r.entity.opts.separator
			}
			e.encode(child, r)
		}
		return
	}
	value := field.ToString()
	if value == RESERVE_LEVEL_PLACEHOLDER {
		value = r.entity.level.String()
	} else if value == RESERVE_MESSAGE_PLACEHOLDER {
		value = *r.msg
	}
	r.output += field.Decorate(value, r.color)
}

func (e *plainEncoder) Encode(r *record) string {
	for index, field := range r.entity.opts.fields {
		if index != 0 {
			r.output += r.entity.opts.separator
		}

		e.encode(field, r)
	}
	return r.output
}

type jsonEncoder struct{}

func (e *jsonEncoder) encode(field LogField, r *record) {

	if len(field.Children()) > 0 {
		key := field.Key()
		r.output += key + ": {"
		for index, child := range field.Children() {
			if index != 0 {
				r.output += ", "
			}
			e.encode(child, r)
		}
		r.output += "}"
		return
	}
	value := field.ToString()
	if value == RESERVE_LEVEL_PLACEHOLDER {
		value = r.entity.level.String()
	} else if value == RESERVE_MESSAGE_PLACEHOLDER {
		value = *r.msg
	}
	key := field.Key()
	r.output += fmt.Sprintf(`"%s": "%s"`, key, value)
}

func (e *jsonEncoder) Encode(r *record) string {
	r.output = "{"
	for index, field := range r.entity.opts.fields {
		if index != 0 {
			r.output += ", "
		}
		e.encode(field, r)
	}
	r.output += "}"

	return r.output
}

type logFmtEncoder struct{}

func (e *logFmtEncoder) encode(field LogField, key string, r *record) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				r.output += " "
			}
			subKey := child.Key()
			subKey = key + "." + subKey
			e.encode(child, subKey, r)
		}
		return
	}

	value := field.ToString()
	if value == RESERVE_LEVEL_PLACEHOLDER {
		value = r.entity.level.String()
	} else if value == RESERVE_MESSAGE_PLACEHOLDER {
		value = *r.msg
	}

	r.output += fmt.Sprintf(`%s="%s"`, key, value)
}

func (e *logFmtEncoder) Encode(r *record) string {
	for index, field := range r.entity.opts.fields {
		if index != 0 {
			r.output += " "
		}
		key := field.Key()
		e.encode(field, key, r)
	}
	return r.output
}
