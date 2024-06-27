package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	RESERVE_MESSAGE_PLACEHOLDER = "__MESSAGE__"
	RESERVE_LEVEL_PLACEHOLDER   = "__LEVEL__"
)

type LogField interface {
	Key() string
	ToString() string
	Children() []LogField
	Decorate(s string, color bool) string

	setKey(key string)
	setColor(color Color)
	setBackgroundColor(color Color)
	setPrefix(prefix string)
	setSuffix(suffix string)
	setUpper()
	setLower()
}

var _ LogField = (*Field)(nil)

type Field struct {
	key        string
	color      Color
	background Color
	prefix     string
	suffix     string
	upper      bool
	lower      bool
	children   []LogField
}

func (f *Field) Decorate(s string, color bool) string {
	if f.upper && !f.lower {
		s = strings.ToUpper(s)
	}
	if !f.upper && f.lower {
		s = strings.ToLower(s)
	}
	if color {
		s = f.Color(s)
	}
	s = fmt.Sprintf("%s%s%s", f.prefix, s, f.suffix)
	return s
}

func (f *Field) Color(text string) string {

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

func (f *Field) Key() string {
	return f.key
}

func (f *Field) ToString() string {
	return ""
}

func (f *Field) Children() []LogField {
	return f.children
}

func (f *Field) setKey(key string) {
	f.key = key
}

func (f *Field) setColor(color Color) {
	f.color = color
}

func (f *Field) setBackgroundColor(color Color) {
	f.background = color.ToBackground()
}
func (f *Field) setPrefix(prefix string) {
	f.prefix = prefix
}

func (f *Field) setSuffix(suffix string) {
	f.suffix = suffix
}

func (f *Field) setUpper() {
	f.upper = true
}

func (f *Field) setLower() {
	f.lower = true
}

type FieldBuilder struct {
	field LogField
}

func (b *FieldBuilder) Key(key string) *FieldBuilder {
	b.field.setKey(key)
	return b
}

func (b *FieldBuilder) Color(color Color) *FieldBuilder {
	b.field.setColor(color)
	return b
}

func (b *FieldBuilder) Background(color Color) *FieldBuilder {
	b.field.setBackgroundColor(color)
	return b
}

func (b *FieldBuilder) Prefix(prefix string) *FieldBuilder {
	b.field.setPrefix(prefix)
	return b
}

func (b *FieldBuilder) Suffix(suffix string) *FieldBuilder {
	b.field.setSuffix(suffix)
	return b
}

func (b *FieldBuilder) Upper() *FieldBuilder {
	b.field.setUpper()
	return b
}

func (b *FieldBuilder) Lower() *FieldBuilder {
	b.field.setLower()
	return b
}

func (b *FieldBuilder) Build() LogField {
	return b.field
}

// ====== FastField ======
type fastField struct {
	*Field
	msg interface{}
}

func (f *fastField) ToString() string {
	return fmt.Sprintf("%v", f.msg)
}

func F(key string, value interface{}) *FieldBuilder {
	f := &fastField{
		Field: &Field{key: key},
		msg:   value,
	}
	return &FieldBuilder{field: f}
}

// ====== Group ======
type group struct {
	*Field
}

func Group(key string, fields ...LogField) *FieldBuilder {
	g := &group{
		Field: &Field{
			key:      key,
			children: fields,
		},
	}
	return &FieldBuilder{field: g}
}

// ====== LevelField ======
type levelField struct {
	Field
	upper bool
}

func (f *levelField) ToString() string {
	return RESERVE_LEVEL_PLACEHOLDER
}

func LevelField(upper bool) *FieldBuilder {
	f := &levelField{
		upper: true,
	}
	return &FieldBuilder{field: f}
}

// ====== MessageField ======
type messageField struct {
	Field
}

func (f *messageField) ToString() string {
	return RESERVE_MESSAGE_PLACEHOLDER
}

func MessageField() *FieldBuilder {
	f := &messageField{}
	return &FieldBuilder{field: f}
}

// ====== DatetimeField ======
type datetimeField struct {
	Field
	layout string
}

func (f *datetimeField) ToString() string {
	if f.layout == "" {
		f.layout = "2006-01-02 15:04:05.000"
	}
	s := time.Now().Format(f.layout)
	return s
}

func DatetimeField(layout string) *FieldBuilder {
	f := &datetimeField{
		layout: layout,
	}
	return &FieldBuilder{field: f}
}

// ====== TimeField ======

type UnixTimeType uint8

const (
	Unix = iota
	UnixMilli
	UnixMicro
	UnixNano
)

func (f *timeField) ToString() string {
	s := fmt.Sprintf("%d", time.Now().UnixMilli())
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
	return s
}

type timeField struct {
	Field
	t UnixTimeType
}

func TimeField(t UnixTimeType) *FieldBuilder {
	f := &timeField{t: t}
	return &FieldBuilder{field: f}
}

// ====== FuncNameField ======
type funcNameField struct {
	Field
	shorten bool
}

func (f *funcNameField) ToString() string {
	pc, _, _, ok := runtime.Caller(7)
	if !ok {
		return "???"
	}
	funcName := runtime.FuncForPC(pc).Name()
	if f.shorten {
		name := funcName
		for i := len(name) - 1; i > 0; i-- {
			if funcName[i] == '.' {
				name = funcName[i+1:]
				break
			}
		}
		funcName = name
	}

	return funcName
}

func FuncNameField(shorten bool) *FieldBuilder {
	f := &funcNameField{
		shorten: shorten,
	}
	return &FieldBuilder{field: f}
}

// ====== CallerField ======
type callerField struct {
	Field
	shorten      bool
	showFuncName bool
}

func (f *callerField) ToString() string {
	pc, file, line, ok := runtime.Caller(7)
	if !ok {
		file = "???"
		line = 0
	}
	if f.shorten {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	var funcName string
	if f.showFuncName {
		funcName = runtime.FuncForPC(pc).Name()
		if f.shorten {
			name := funcName
			for i := len(name) - 1; i > 0; i-- {
				if funcName[i] == '.' {
					name = funcName[i+1:]
					break
				}
			}
			funcName = name
		}
		funcName = " " + funcName
	}

	s := fmt.Sprintf("%s %d%s", file, line, funcName)
	return s
}

func CallerField(shorten bool, showFuncName bool) *FieldBuilder {
	f := &callerField{
		shorten:      shorten,
		showFuncName: showFuncName,
	}
	return &FieldBuilder{field: f}
}

// ====== CustomField ======

type customField struct {
	Field
	handler func() string
}

func (f *customField) ToString() string {
	if f.handler == nil {
		return ""
	}
	s := f.handler()
	return s
}

func CustomField(handler func() string) *FieldBuilder {
	f := &customField{
		handler: handler,
	}
	return &FieldBuilder{field: f}
}
