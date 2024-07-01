package logger

import (
	"github.com/yates-z/easel/internal/pool"
	"github.com/yates-z/easel/logger/buffer"
	"runtime"
	"time"
)

type FieldType uint8

const (
	CommonFieldType FieldType = iota
	LevelFieldType
	MessageFieldType
	GroupFieldType
)

type LogField interface {
	Kind() FieldType
	Key() string
	Log(buf *buffer.Buffer)
	Children() []LogField
	Free()
	decorate(buf *buffer.Buffer, s string, withColor bool)

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
	_child     LogField
}

func (f *Field) Kind() FieldType {
	return CommonFieldType
}

func (f *Field) decorate(buf *buffer.Buffer, s string, withColor bool) {

	_, _ = buf.WriteString(f.prefix)
	if withColor && !(f.color.IsDefault() && f.background.IsDefault()) {
		_ = buf.WriteByte(0x1B)
		if f.color.IsDefault() && !f.background.IsDefault() {
			_, _ = buf.WriteString("[0;0;" + f.background.String() + "m")
		} else {
			_, _ = buf.WriteString("[0;" + f.background.String() + ";" + f.color.String() + "m")
		}
	}

	if len(s) == 0 {
		f._child.Log(buf)
	} else {
		_, _ = buf.WriteString(s)

		if f.upper && !f.lower {
			for i := 1; i <= len(s); i++ {
				if (*buf)[buf.Len()-i] >= 'a' && (*buf)[buf.Len()-i] <= 'z' {
					(*buf)[buf.Len()-i] -= 'a' - 'A'
				}
			}
		}
		if !f.upper && f.lower {
			for i := 1; i <= len(s); i++ {
				if (*buf)[buf.Len()-i] >= 'A' && (*buf)[buf.Len()-i] <= 'Z' {
					(*buf)[buf.Len()-i] += 'a' - 'A'
				}
			}
		}
	}
	if withColor && !(f.color.IsDefault() && f.background.IsDefault()) {
		_ = buf.WriteByte(0x1B)
		_, _ = buf.WriteString("[0m")
	}
	_, _ = buf.WriteString(f.suffix)

}

func (f *Field) Key() string {
	return f.key
}

func (f *Field) Log(_ *buffer.Buffer) {}

func (f *Field) Children() []LogField {
	return f.children
}

func (f *Field) Free() {

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

func (b FieldBuilder) Key(key string) FieldBuilder {
	b.field.setKey(key)
	return b
}

func (b FieldBuilder) Color(color Color) FieldBuilder {
	b.field.setColor(color)
	return b
}

func (b FieldBuilder) Background(color Color) FieldBuilder {
	b.field.setBackgroundColor(color)
	return b
}

func (b FieldBuilder) Prefix(prefix string) FieldBuilder {
	b.field.setPrefix(prefix)
	return b
}

func (b FieldBuilder) Suffix(suffix string) FieldBuilder {
	b.field.setSuffix(suffix)
	return b
}

func (b FieldBuilder) Upper() FieldBuilder {
	b.field.setUpper()
	return b
}

func (b FieldBuilder) Lower() FieldBuilder {
	b.field.setLower()
	return b
}

func (b FieldBuilder) Build() LogField {
	return b.field
}

// ====== FastField ======
type fastField struct {
	*Field
	msg []byte
}

var fastFieldPool = pool.New(func() *fastField {
	f := &fastField{
		Field: &Field{},
		msg:   []byte{},
	}
	f.Field._child = f
	return f
})

func (f *fastField) Free() {
	f.msg = f.msg[:0]
	f.color = DefaultColor
	f.background = DefaultColor
	f.prefix = ""
	f.suffix = ""
	f.upper = false
	f.lower = false
	fastFieldPool.Put(f)
}

func (f *fastField) Color(color Color) *fastField {
	f.color = color
	return f
}

func (f *fastField) Log(buf *buffer.Buffer) {
	_, _ = buf.Write(f.msg)
}

func F(key string, value string) FieldBuilder {
	f := fastFieldPool.Get()
	f.Field.setKey(key)
	f.msg = append(f.msg, value...)
	return FieldBuilder{field: f}
}

// ====== Group ======
type group struct {
	*Field
}

var groupPool = pool.New(func() *group {
	g := &group{
		Field: &Field{},
	}
	g.Field._child = g
	return g
})

func (f *group) Free() {
	for _, child := range f.Field.children {
		child.Free()
	}
	f.children = f.children[:0]
	groupPool.Put(f)
}

func (f *group) Kind() FieldType {
	return GroupFieldType
}

func Group(key string, fields ...FieldBuilder) FieldBuilder {
	g := groupPool.Get()
	g.Field.setKey(key)
	for _, f := range fields {
		g.Field.children = append(g.Field.children, f.Build())
	}
	return FieldBuilder{field: g}
}

// ====== LevelField ======
type levelField struct {
	Field
}

func (f *levelField) Kind() FieldType {
	return LevelFieldType
}

func LevelField() FieldBuilder {
	f := &levelField{Field{key: "level"}}
	f.Field._child = f
	return FieldBuilder{field: f}
}

// ====== MessageField ======
type messageField struct {
	Field
}

func (f *messageField) Kind() FieldType {
	return MessageFieldType
}

func MessageField() FieldBuilder {
	f := &messageField{Field{key: "msg"}}
	f.Field._child = f
	return FieldBuilder{field: f}
}

// ====== DatetimeField ======
type datetimeField struct {
	Field
	layout string
}

func (f *datetimeField) Log(buf *buffer.Buffer) {
	now := time.Now()
	if f.layout == "" {
		f.layout = "2006/01/02 15:04:05.000"
	}
	*buf = now.AppendFormat(*buf, f.layout)
}

func DatetimeField(layout string) FieldBuilder {
	f := &datetimeField{
		Field:  Field{key: "datetime"},
		layout: layout,
	}
	f.Field._child = f
	return FieldBuilder{field: f}
}

// ====== TimeField ======

type UnixTimeType uint8

const (
	Unix = iota
	UnixMilli
	UnixMicro
	UnixNano
)

func (f *timeField) Log(buf *buffer.Buffer) {
	switch f.t {
	case Unix:
		buf.WriteInt(time.Now().Unix())
	case UnixMilli:
		buf.WriteInt(time.Now().UnixMilli())
	case UnixMicro:
		buf.WriteInt(time.Now().UnixMicro())
	case UnixNano:
		buf.WriteInt(time.Now().UnixNano())
	}
}

type timeField struct {
	Field
	t UnixTimeType
}

func TimeField(t UnixTimeType) FieldBuilder {
	f := &timeField{
		Field: Field{key: "time"},
		t:     t,
	}
	f.Field._child = f
	return FieldBuilder{field: f}
}

// ====== FuncNameField ======
type funcNameField struct {
	Field
	shorten bool
}

func (f *funcNameField) Log(buf *buffer.Buffer) {
	var pcs [1]uintptr
	runtime.Callers(8, pcs[:])
	if pcs[0] == 0 {
		_, _ = buf.WriteString("???")
	}
	funcName := runtime.FuncForPC(pcs[0]).Name()

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
	_, _ = buf.WriteString(funcName)
}

func FuncNameField(shorten bool) FieldBuilder {
	f := &funcNameField{
		Field:   Field{key: "func"},
		shorten: shorten,
	}
	f.Field._child = f
	return FieldBuilder{field: f}
}

// ====== CallerField ======
type callerField struct {
	Field
	shorten      bool
	showFuncName bool
}

func (f *callerField) Log(buf *buffer.Buffer) {
	pcs := make([]uintptr, 1)
	runtime.Callers(7, pcs[:])
	if pcs[0] == 0 {
		_, _ = buf.WriteString("??? 0")
	}
	fs := runtime.CallersFrames(pcs)
	fs.Next()
	//if f.shorten {
	//	short := file
	//	for i := len(file) - 1; i > 0; i-- {
	//		if file[i] == '/' {
	//			short = file[i+1:]
	//			break
	//		}
	//	}
	//	file = short
	//}
	//
	//var funcName string
	//if f.showFuncName {
	//	funcName = runtime.FuncForPC(pc).Name()
	//	if f.shorten {
	//		name := funcName
	//		for i := len(name) - 1; i > 0; i-- {
	//			if funcName[i] == '.' {
	//				name = funcName[i+1:]
	//				break
	//			}
	//		}
	//		funcName = name
	//	}
	//	funcName = " " + funcName
	//}
	//
	//_ = fmt.Sprintf("%s %d%s", file, line, funcName)

}

func CallerField(shorten bool, showFuncName bool) FieldBuilder {
	f := &callerField{
		Field:        Field{key: "caller"},
		shorten:      shorten,
		showFuncName: showFuncName,
	}
	f.Field._child = f
	return FieldBuilder{field: f}
}

// ====== CustomField ======

type customField struct {
	Field
	handler func(buf *buffer.Buffer)
}

func (f *customField) Log(buf *buffer.Buffer) {
	if f.handler == nil {
		return
	}
	f.handler(buf)
}

func CustomField(handler func(buf *buffer.Buffer)) FieldBuilder {
	f := &customField{
		handler: handler,
	}
	f.Field._child = f
	return FieldBuilder{field: f}
}
