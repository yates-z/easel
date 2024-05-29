package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type LogField interface {
	Key() (s string, ps string)
	ToString() (s string, ps string)
	Text(*string) (s string, ps string)
	paintString(string) string
}

type FieldConstructor = func(entity *logEntity) LogField

type Field struct {
	key        string
	entity     *logEntity
	color      Color
	background Color
}

func (f *Field) Key() (s string, ps string) {
	s = f.key
	ps = f.paintString(s)
	return
}

func (f *Field) Text(_ *string) (s string, ps string) {
	return
}

func (f *Field) paintString(text string) string {

	if !f.color.IsDefault() && f.background.IsDefault() {
		// only foreground color
		text = f.color.Paint(text)

	} else if f.color.IsDefault() && !f.background.IsDefault() {
		// only background color
		text = f.background.Paint(text)

	} else if !f.color.IsDefault() && !f.background.IsDefault() {
		// both foreground and background
		text = f.color.PaintWith(f.background, text)

	} else {
		// neither
	}
	return text
}

// ====== LevelField ======
type levelField struct {
	Field
	upper bool
}

func (f *levelField) ToString() (s string, ps string) {
	s = f.entity.level.String()
	if f.upper {
		s = strings.ToUpper(s)
	}

	ps = f.paintString(s)
	return
}

type levelFieldBuilder struct {
	field levelField
}

func (f *levelFieldBuilder) Upper(upper bool) *levelFieldBuilder {
	f.field.upper = upper
	return f
}

func (f *levelFieldBuilder) Color(color Color) *levelFieldBuilder {
	f.field.color = color
	return f
}

func (f *levelFieldBuilder) Background(color Color) *levelFieldBuilder {
	f.field.background = color.ToBackground()
	return f
}

func (f *levelFieldBuilder) Build() FieldConstructor {

	return func(entity *logEntity) LogField {
		field := f.field
		field.entity = entity
		return &field
	}
}

func LevelField(key string) *levelFieldBuilder {
	f := levelField{
		Field: Field{key: key},
	}
	return &levelFieldBuilder{field: f}
}

//func LevelField(key string, upper bool) LogFieldHandler {
//	f := levelField{
//		Field: Field{Key: key},
//		upper: upper,
//	}
//	return func(entity *logEntity) LogField {
//		field := f
//		field.entity = entity
//		return &field
//	}
//}

// ====== DatetimeField ======
type datetimeField struct {
	Field
	layout string
}

func (f *datetimeField) ToString() (s string, ps string) {
	if f.layout == "" {
		f.layout = "2006-01-02 15:04:05.000"
	}
	s = time.Now().Format(f.layout)
	ps = f.paintString(s)
	return
}

type datetimeFieldBuilder struct {
	field datetimeField
}

func (f *datetimeFieldBuilder) Layout(layout string) *datetimeFieldBuilder {
	f.field.layout = layout
	return f
}

func (f *datetimeFieldBuilder) Color(color Color) *datetimeFieldBuilder {
	f.field.color = color
	return f
}

func (f *datetimeFieldBuilder) Background(color Color) *datetimeFieldBuilder {
	f.field.background = color.ToBackground()
	return f
}

func (f *datetimeFieldBuilder) Build() FieldConstructor {

	return func(entity *logEntity) LogField {
		field := f.field
		field.entity = entity
		return &field
	}
}

func DatetimeField(key string) *datetimeFieldBuilder {
	f := datetimeField{
		Field: Field{key: key},
	}
	return &datetimeFieldBuilder{field: f}
}

// ====== TimeField ======

type UnixTimeType uint8

const (
	Unix = iota
	UnixMilli
	UnixMicro
	UnixNano
)

func (f *timeField) ToString() (s string, ps string) {
	s = fmt.Sprintf("%d", time.Now().UnixMilli())
	switch f.t {
	case Unix:
		s = fmt.Sprintf("%d", time.Now().Unix())
	case UnixMilli:
		s = fmt.Sprintf("%d", time.Now().UnixMilli())
	case UnixMicro:
		s = fmt.Sprintf("%d", time.Now().UnixMicro())
	case UnixNano:
		s = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	ps = f.paintString(s)
	return
}

type timeField struct {
	Field
	t UnixTimeType
}

type timeFieldBuilder struct {
	field timeField
}

func (f *timeFieldBuilder) Type(t UnixTimeType) *timeFieldBuilder {
	f.field.t = t
	return f
}

func (f *timeFieldBuilder) Color(color Color) *timeFieldBuilder {
	f.field.color = color
	return f
}

func (f *timeFieldBuilder) Background(color Color) *timeFieldBuilder {
	f.field.background = color.ToBackground()
	return f
}

func (f *timeFieldBuilder) Build() FieldConstructor {

	return func(entity *logEntity) LogField {
		field := f.field
		field.entity = entity
		return &field
	}
}

func TimeField(key string) *timeFieldBuilder {
	f := timeField{
		Field: Field{key: key},
	}
	return &timeFieldBuilder{field: f}
}

// ====== BodyField ======
type bodyField struct {
	Field
}

func (f *bodyField) ToString() (s string, ps string) {
	return
}

func (f *bodyField) Text(text *string) (s string, ps string) {
	s = *text
	ps = f.paintString(s)
	return
}

type bodyFieldBuilder struct {
	field bodyField
}

func (f *bodyFieldBuilder) Color(color Color) *bodyFieldBuilder {
	f.field.color = color
	return f
}

func (f *bodyFieldBuilder) Background(color Color) *bodyFieldBuilder {
	f.field.background = color.ToBackground()
	return f
}

func (f *bodyFieldBuilder) Build() FieldConstructor {

	return func(entity *logEntity) LogField {
		field := f.field
		field.entity = entity
		return &field
	}
}

func BodyField(key string) *bodyFieldBuilder {
	f := bodyField{
		Field: Field{key: key},
	}
	return &bodyFieldBuilder{field: f}
}

// ====== LongFileField ======
type longFileField struct {
	Field
}

func (f *longFileField) ToString() (s string, ps string) {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "???"
		line = 0
	}
	s = fmt.Sprintf("%s %d", file, line)
	ps = f.paintString(s)
	return
}

type longFileFieldBuilder struct {
	field longFileField
}

func (f *longFileFieldBuilder) Color(color Color) *longFileFieldBuilder {
	f.field.color = color
	return f
}

func (f *longFileFieldBuilder) Background(color Color) *longFileFieldBuilder {
	f.field.background = color.ToBackground()
	return f
}

func (f *longFileFieldBuilder) Build() FieldConstructor {

	return func(entity *logEntity) LogField {
		field := f.field
		field.entity = entity
		return &field
	}
}

func LongFileField(key string) *longFileFieldBuilder {
	f := longFileField{
		Field: Field{key: key},
	}
	return &longFileFieldBuilder{field: f}
}

// ====== ShortFileField ======
type shortFileField struct {
	Field
}

func (f *shortFileField) ToString() (s string, ps string) {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	s = fmt.Sprintf("%s %d", file, line)
	ps = f.paintString(s)
	return
}

type shortFileFieldBuilder struct {
	field shortFileField
}

func (f *shortFileFieldBuilder) Color(color Color) *shortFileFieldBuilder {
	f.field.color = color
	return f
}

func (f *shortFileFieldBuilder) Background(color Color) *shortFileFieldBuilder {
	f.field.background = color.ToBackground()
	return f
}

func (f *shortFileFieldBuilder) Build() FieldConstructor {

	return func(entity *logEntity) LogField {
		field := f.field
		field.entity = entity
		return &field
	}
}

func ShortFileField(key string) *shortFileFieldBuilder {
	f := shortFileField{
		Field: Field{key: key},
	}
	return &shortFileFieldBuilder{field: f}
}
