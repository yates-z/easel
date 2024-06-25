package logger

import "fmt"

type Encoder interface {
	Encode(text string) (string, string)
}

type EncoderConstructor = func(entity *logEntity) Encoder

type BaseEncoder struct {
	entity *logEntity
}

type plainEncoder struct {
	BaseEncoder
}

func (e *plainEncoder) encodeField(field LogField, text, output, pOutput *string) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				*output += e.entity.separator
				*pOutput += e.entity.separator
			}
			e.encodeField(child, text, output, pOutput)
		}
		return
	}

	s, ps := field.Text(text)
	*output += s
	*pOutput += ps

	s, ps = field.ToString()
	*output += s
	*pOutput += ps
}

func (e *plainEncoder) Encode(text string) (string, string) {
	output := ""
	// painted output
	pOutput := ""
	for index, f := range e.entity.fields {
		if index != 0 {
			output += e.entity.separator
			pOutput += e.entity.separator
		}

		e.encodeField(f, &text, &output, &pOutput)
	}
	for count := e.entity.skipLines; count >= 0; count-- {
		output += "\n"
		pOutput += "\n"
	}
	return output, pOutput
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

func (e *jsonEncoder) encodeField(field LogField, text, output, pOutput *string) {
	if len(field.Children()) > 0 {
		*output += field.Key() + ": {"
		*pOutput += field.ColoredKey() + ": {"

		for index, child := range field.Children() {
			if index != 0 {
				*output += ", "
				*pOutput += ", "
			}
			e.encodeField(child, text, output, pOutput)
		}
		*output += "}"
		*pOutput += "}"
		return
	}

	value, pValue := field.Text(text)
	s, ps := field.ToString()
	value += s
	pValue += ps

	*output += fmt.Sprintf(`"%s": "%s"`, field.Key(), value)
	*pOutput += fmt.Sprintf(`"%s": "%s"`, field.ColoredKey(), pValue)
}

func (e *jsonEncoder) Encode(text string) (string, string) {
	output := "{"
	// painted output
	pOutput := "{"
	for index, f := range e.entity.fields {
		if index != 0 {
			output += ", "
			pOutput += ", "
		}
		e.encodeField(f, &text, &output, &pOutput)
	}
	output += "}"
	pOutput += "}"

	for count := e.entity.skipLines; count >= 0; count-- {
		output += "\n"
		pOutput += "\n"
	}
	return output, pOutput
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

func (e *logFmtEncoder) encodeField(field LogField, key, coloredKey string, text, output, pOutput *string) {

	if len(field.Children()) > 0 {
		for index, child := range field.Children() {
			if index != 0 {
				*output += " "
				*pOutput += " "
			}
			subKey := child.Key()
			subColoredKey := child.ColoredKey()
			subKey = key + "." + subKey
			subColoredKey = coloredKey + "." + subColoredKey
			e.encodeField(child, subKey, subColoredKey, text, output, pOutput)
		}
		return
	}

	value, pValue := field.Text(text)
	s, ps := field.ToString()
	value += s
	pValue += ps

	*output += fmt.Sprintf(`%s=%s`, key, value)
	*pOutput += fmt.Sprintf(`%s=%s`, coloredKey, pValue)
}

func (e *logFmtEncoder) Encode(text string) (string, string) {
	output, pOutput := "", ""
	for index, f := range e.entity.fields {
		if index != 0 {
			output += " "
			pOutput += " "
		}
		e.encodeField(f, f.Key(), f.ColoredKey(), &text, &output, &pOutput)
	}

	for count := e.entity.skipLines; count >= 0; count-- {
		output += "\n"
		pOutput += "\n"
	}
	return output, pOutput
}

func LogFmtEncoder() EncoderConstructor {
	return func(entity *logEntity) Encoder {
		encoder := logFmtEncoder{}
		encoder.entity = entity
		return &encoder
	}
}
