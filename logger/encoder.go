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

func (e *plainEncoder) Encode(text string) (string, string) {
	output := ""
	// painted output
	pOutput := ""
	for index, f := range e.entity.fields {
		if index != 0 {
			output += e.entity.separator
			pOutput += e.entity.separator
		}

		s, ps := f.Text(&text)
		output += s
		pOutput += ps

		s, ps = f.ToString()
		output += s
		pOutput += ps
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

func (e *jsonEncoder) Encode(text string) (string, string) {
	output := "{"
	// painted output
	pOutput := "{"
	for index, f := range e.entity.fields {
		if index != 0 {
			output += ", "
			pOutput += ", "
		}

		value, pValue := f.Text(&text)
		s, ps := f.ToString()
		value += s
		pValue += ps

		key, pKey := f.Key()
		output += fmt.Sprintf(`"%s": "%s"`, key, value)
		pOutput += fmt.Sprintf(`"%s": "%s"`, pKey, pValue)
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

func (e *logFmtEncoder) Encode(text string) (string, string) {
	output, pOutput := "", ""
	for index, f := range e.entity.fields {
		if index != 0 {
			output += " "
			pOutput += " "
		}

		value, pValue := f.Text(&text)
		s, ps := f.ToString()
		value += s
		pValue += ps

		key, pKey := f.Key()
		output += fmt.Sprintf(`%s=%s`, key, value)
		pOutput += fmt.Sprintf(`%s=%s`, pKey, pValue)
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
