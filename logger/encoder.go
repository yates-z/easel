package logger

import "fmt"

type Encoder interface {
	Encode(text string) string
	EncodeWithColor(text string) string
}

type EncoderConstructor = func(entity *logEntity) Encoder

type BaseEncoder struct {
	entity *logEntity
}

type plainEncoder struct {
	BaseEncoder
}

func (e *plainEncoder) encodeField(field LogField, text, output *string, color bool) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				*output += e.entity.separator
			}
			e.encodeField(child, text, output, color)
		}
		return
	}

	s := field.Text(text)
	if color {
		s = field.Color(s)
	}
	*output += s

	s = field.ToString()
	if color {
		s = field.Color(s)
	}
	*output += s
}

func (e *plainEncoder) encode(text string, color bool) string {
	output := ""
	for index, f := range e.entity.fields {
		if index != 0 {
			output += e.entity.separator
		}

		e.encodeField(f, &text, &output, color)
	}
	for count := e.entity.skipLines; count >= 0; count-- {
		output += "\n"
	}
	return output
}

func (e *plainEncoder) Encode(text string) string {
	return e.encode(text, false)
}

func (e *plainEncoder) EncodeWithColor(text string) string {
	return e.encode(text, true)
}

func PlainEncoder() EncoderConstructor {
	return func(entity *logEntity) Encoder {
		encoder := plainEncoder{}
		encoder.entity = entity
		return &encoder
	}
}

type jsonEncoder struct {
	BaseEncoder
}

func (e *jsonEncoder) encodeField(field LogField, text, output *string, color bool) {
	if len(field.Children()) > 0 {
		key := field.Key()
		if color {
			key = field.Color(key)
		}
		*output += key + ": {"
		for index, child := range field.Children() {
			if index != 0 {
				*output += ", "
			}
			e.encodeField(child, text, output, color)
		}
		*output += "}"
		return
	}
	value := field.Text(text)
	s := field.ToString()
	value += s

	key := field.Key()
	if color {
		value = field.Color(value)
		key = field.Color(key)
	}
	*output += fmt.Sprintf(`"%s": "%s"`, key, value)

}

func (e *jsonEncoder) encode(text string, color bool) string {
	output := "{"
	for index, f := range e.entity.fields {
		if index != 0 {
			output += ", "
		}
		e.encodeField(f, &text, &output, color)
	}
	output += "}"

	for count := e.entity.skipLines; count >= 0; count-- {
		output += "\n"
	}
	return output
}

func (e *jsonEncoder) Encode(text string) string {
	return e.encode(text, false)
}

func (e *jsonEncoder) EncodeWithColor(text string) string {
	return e.encode(text, true)
}

func JsonEncoder() EncoderConstructor {
	return func(entity *logEntity) Encoder {
		encoder := jsonEncoder{}
		encoder.entity = entity
		return &encoder
	}
}

type logFmtEncoder struct {
	BaseEncoder
}

func (e *logFmtEncoder) encodeField(field LogField, key string, text, output *string, color bool) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				*output += " "
			}
			subKey := child.Key()
			subKey = key + "." + subKey
			e.encodeField(child, subKey, text, output, color)
		}
		return
	}

	value := field.Text(text)
	s := field.ToString()
	value += s
	if color {
		value = field.Color(value)
		key = field.Color(key)
	}

	*output += fmt.Sprintf(`%s=%s`, key, value)
}

func (e *logFmtEncoder) encode(text string, color bool) string {
	output := ""
	for index, f := range e.entity.fields {
		if index != 0 {
			output += " "
		}
		key := f.Key()
		if color {
			key = f.Color(key)
		}
		e.encodeField(f, key, &text, &output, color)
	}

	for count := e.entity.skipLines; count >= 0; count-- {
		output += "\n"
	}
	return output
}

func (e *logFmtEncoder) Encode(text string) string {
	return e.encode(text, false)
}

func (e *logFmtEncoder) EncodeWithColor(text string) string {
	return e.encode(text, true)
}

func LogFmtEncoder() EncoderConstructor {
	return func(entity *logEntity) Encoder {
		encoder := logFmtEncoder{}
		encoder.entity = entity
		return &encoder
	}
}
