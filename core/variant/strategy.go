package variant

import (
	"encoding/binary"
	"math"
	"strconv"
	"time"
	"unsafe"
)

type IConvertStrategy[T any] interface {
	FromString(v Variant) T
	FromBool(v Variant) T
	FromInt(v Variant) T
	FromInt8(v Variant) T
	FromInt16(v Variant) T
	FromInt32(v Variant) T
	FromInt64(v Variant) T
	FromUint(v Variant) T
	FromUint8(v Variant) T
	FromUint16(v Variant) T
	FromUint32(v Variant) T
	FromUint64(v Variant) T
	FromFloat32(v Variant) T
	FromFloat64(v Variant) T
	FromTime(v Variant) T
}

var Strategies = strategies{
	string:  newStringConverter(),
	int:     newIntConverter(),
	uint:    newUintConverter(),
	float32: newFloat32Converter(),
	float64: newFloat64Converter(),
	time:    newTimeConverter(),
}

type strategies struct {
	string  *stringConverter
	int     *intConverter
	uint    *uintConverter
	float32 *float32Converter
	float64 *float64Converter
	time    *timeConverter
}

type Converter[T any] struct {
	m map[Kind]func(v Variant) T
}

func (c Converter[T]) Get(k Kind) func(v Variant) T {
	return c.m[k]
}

var _ IConvertStrategy[string] = (*stringConverter)(nil)

type stringConverter struct {
	Converter[string]
}

func (c stringConverter) FromString(v Variant) string {
	return *(*string)(unsafe.Pointer(&v.Data))
}

func (c stringConverter) FromBool(v Variant) string {
	if len(v.Data) == 0 || v.Data[0] == 0x00 {
		return "false"
	}
	return "true"
}

func (c stringConverter) FromInt(v Variant) string {
	var i int
	if len(v.Data) == 4 && intSize == 32 {
		i = int(binary.BigEndian.Uint32(v.Data))
	} else if len(v.Data) == 8 && intSize == 64 {
		i = int(binary.BigEndian.Uint64(v.Data))
	}
	return strconv.Itoa(i)
}

func (c stringConverter) FromInt8(v Variant) string {
	var i int8
	if len(v.Data) == 1 {
		i = int8(v.Data[0])
	}
	b := strconv.AppendInt(make([]byte, 0), int64(i), 10)
	return *(*string)(unsafe.Pointer(&b))
}

func (c stringConverter) FromInt16(v Variant) string {
	var i int16
	if len(v.Data) == 2 {
		i = int16(binary.BigEndian.Uint16(v.Data))
	}
	return strconv.Itoa(int(i))
}

func (c stringConverter) FromInt32(v Variant) string {
	i := int32(binary.BigEndian.Uint32(v.Data))
	return strconv.Itoa(int(i))
}

func (c stringConverter) FromInt64(v Variant) string {
	i := int64(binary.BigEndian.Uint64(v.Data))
	return strconv.Itoa(int(i))
}

func (c stringConverter) FromUint(v Variant) string {
	var i uint64
	if intSize == 32 {
		i = uint64(binary.BigEndian.Uint32(v.Data))
	} else if intSize == 64 {
		i = binary.BigEndian.Uint64(v.Data)
	}
	return strconv.FormatUint(i, 10)
}

func (c stringConverter) FromUint8(v Variant) string {
	i := v.Data[0]
	b := strconv.AppendUint(make([]byte, 0), uint64(i), 10)
	return *(*string)(unsafe.Pointer(&b))
}

func (c stringConverter) FromUint16(v Variant) string {
	i := uint64(binary.BigEndian.Uint16(v.Data))
	return strconv.FormatUint(i, 10)
}

func (c stringConverter) FromUint32(v Variant) string {
	i := uint64(binary.BigEndian.Uint32(v.Data))
	return strconv.FormatUint(i, 10)
}

func (c stringConverter) FromUint64(v Variant) string {
	i := binary.BigEndian.Uint64(v.Data)
	return strconv.FormatUint(i, 10)
}

func (c stringConverter) FromFloat32(v Variant) string {
	var f float32
	if len(v.Data) == 4 {
		f = math.Float32frombits(binary.BigEndian.Uint32(v.Data))
	}
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

func (c stringConverter) FromFloat64(v Variant) string {
	var f float64
	if len(v.Data) == 8 {
		f = math.Float64frombits(binary.BigEndian.Uint64(v.Data))
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (c stringConverter) FromTime(v Variant) string {
	var t time.Time
	err := t.UnmarshalBinary(v.Data)
	if err != nil {
		return ""
	}
	return t.Format(v.layout)
}

func newStringConverter() *stringConverter {
	c := &stringConverter{}
	c.m = map[Kind]func(v Variant) string{
		String:  c.FromString,
		Bool:    c.FromBool,
		Int:     c.FromInt,
		Int8:    c.FromInt8,
		Int16:   c.FromInt16,
		Int32:   c.FromInt32,
		Int64:   c.FromInt64,
		Uint:    c.FromUint,
		Uint8:   c.FromUint8,
		Uint16:  c.FromUint16,
		Uint32:  c.FromUint32,
		Uint64:  c.FromUint64,
		Float32: c.FromFloat32,
		Float64: c.FromFloat64,
		Time:    c.FromTime,
	}
	return c
}

var _ IConvertStrategy[int] = (*intConverter)(nil)

type intConverter struct {
	Converter[int]
}

func (c intConverter) FromString(v Variant) int {

	for index, ch := range v.Data {
		if ch == '.' {
			v.Data = v.Data[:index]
			break
		}
	}

	start := 0
	sLen := len(v.Data)
	if sLen == 0 {
		return 0
	}
	if v.Data[0] == '-' || v.Data[0] == '+' {
		start = 1
		sLen -= 1
		if len(v.Data) < 2 {
			return 0
		}
	}

	n := 0

	if (intSize == 32 && sLen < 10) || (intSize == 64 && sLen < 19) {
		for _, ch := range v.Data[start:] {
			ch -= '0'
			if ch > 9 {
				return 0
			}
			n = n*10 + int(ch)
		}
	}

	sub := 0
	if (intSize == 32 && sLen == 10) || (intSize == 64 && sLen == 19) {
		for _, ch := range v.Data[start : len(v.Data)-1] {
			ch -= '0'
			if ch > 9 {
				return 0
			}
			n = n*10 + int(ch)
		}

		cutoff := maxInt / 10
		remainder := maxInt % 10
		ch := v.Data[len(v.Data)-1] - '0'
		if ch > 9 {
			return 0
		}
		digit := int(ch)
		if v.Data[0] == '-' && cutoff == n && remainder == digit-1 {
			sub = 1
			n = n*10 + remainder
		} else if (cutoff == n && remainder >= digit) || cutoff > n {
			n = n*10 + digit
		} else {
			return 0
		}
	}

	if v.Data[0] == '-' {
		n = -n - sub
	}
	return n
}

func (c intConverter) FromBool(v Variant) int {
	if len(v.Data) == 0 || v.Data[0] == 0x00 {
		return 0
	}
	return 1
}

func (c intConverter) FromInt(v Variant) int {
	var i int
	if len(v.Data) == 4 && intSize == 32 {
		i = int(binary.BigEndian.Uint32(v.Data))
	} else if len(v.Data) == 8 && intSize == 64 {
		i = int(binary.BigEndian.Uint64(v.Data))
	}
	return i
}

func (c intConverter) FromInt8(v Variant) int {
	var i int8
	if len(v.Data) == 1 {
		i = int8(v.Data[0])
	}
	return int(i)
}

func (c intConverter) FromInt16(v Variant) int {
	var i int16
	if len(v.Data) == 2 {
		i = int16(binary.BigEndian.Uint16(v.Data))
	}
	return int(i)
}

func (c intConverter) FromInt32(v Variant) int {
	return int(int32(binary.BigEndian.Uint32(v.Data)))
}

func (c intConverter) FromInt64(v Variant) int {
	return int(int64(binary.BigEndian.Uint64(v.Data)))
}

func (c intConverter) FromUint(v Variant) int {
	var i uint
	if intSize == 32 {
		i = uint(binary.BigEndian.Uint32(v.Data))
		if i > 1<<31-1 {
			return 0
		}
	} else if intSize == 64 {
		i = uint(binary.BigEndian.Uint64(v.Data))
		if i > 1<<63-1 {
			return 0
		}
	}
	return int(i)
}

func (c intConverter) FromUint8(v Variant) int {
	return int(v.Data[0])
}

func (c intConverter) FromUint16(v Variant) int {
	return int(binary.BigEndian.Uint16(v.Data))
}

func (c intConverter) FromUint32(v Variant) int {
	return int(binary.BigEndian.Uint32(v.Data))
}

func (c intConverter) FromUint64(v Variant) int {
	return int(binary.BigEndian.Uint64(v.Data))
}

func (c intConverter) FromFloat32(v Variant) int {
	var f float32
	if len(v.Data) == 4 {
		f = math.Float32frombits(binary.BigEndian.Uint32(v.Data))
	}
	return int(f)
}

func (c intConverter) FromFloat64(v Variant) int {
	var f float64
	if len(v.Data) == 8 {
		f = math.Float64frombits(binary.BigEndian.Uint64(v.Data))
	}
	return int(f)
}

func (c intConverter) FromTime(v Variant) int {
	var t time.Time
	err := t.UnmarshalBinary(v.Data)
	if err != nil {
		return 0
	}
	return int(t.UnixNano())
}

func newIntConverter() *intConverter {
	c := &intConverter{}
	c.m = map[Kind]func(v Variant) int{
		String:  c.FromString,
		Bool:    c.FromBool,
		Int:     c.FromInt,
		Int8:    c.FromInt8,
		Int16:   c.FromInt16,
		Int32:   c.FromInt32,
		Int64:   c.FromInt64,
		Uint:    c.FromUint,
		Uint8:   c.FromUint8,
		Uint16:  c.FromUint16,
		Uint32:  c.FromUint32,
		Uint64:  c.FromUint64,
		Float32: c.FromFloat32,
		Float64: c.FromFloat64,
		Time:    c.FromTime,
	}
	return c
}

var _ IConvertStrategy[uint] = (*uintConverter)(nil)

type uintConverter struct {
	Converter[uint]
}

func (u uintConverter) FromString(v Variant) uint {
	for index, ch := range v.Data {
		if ch == '.' {
			v.Data = v.Data[:index]
			break
		}
	}
	start := 0
	sLen := len(v.Data)
	if sLen == 0 || v.Data[0] == '-' {
		return 0
	}
	if v.Data[0] == '+' {
		start = 1
		sLen -= 1
		if len(v.Data) < 2 {
			return 0
		}
	}

	var n uint = 0

	if (intSize == 32 && sLen < 10) || (intSize == 64 && sLen < 20) {
		for _, ch := range v.Data[start:] {
			ch -= '0'
			if ch > 9 {
				return 0
			}
			n = n*10 + uint(ch)
		}
	}

	if (intSize == 32 && sLen == 10) || (intSize == 64 && sLen == 20) {
		for _, ch := range v.Data[start : len(v.Data)-1] {
			ch -= '0'
			if ch > 9 {
				return 0
			}
			n = n*10 + uint(ch)
		}

		var cutoff uint = maxUint / 10
		var remainder uint = maxUint % 10
		ch := v.Data[len(v.Data)-1] - '0'
		if ch > 9 {
			return 0
		}
		digit := uint(ch)
		if (cutoff == n && remainder >= digit) || cutoff > n {
			n = n*10 + digit
		} else {
			return 0
		}
	}
	return n
}

func (u uintConverter) FromBool(v Variant) uint {
	if len(v.Data) == 0 || v.Data[0] == 0x00 {
		return 0
	}
	return 1
}

func (u uintConverter) FromInt(v Variant) uint {
	var i int
	if len(v.Data) == 4 && intSize == 32 {
		i = int(binary.BigEndian.Uint32(v.Data))
	} else if len(v.Data) == 8 && intSize == 64 {
		i = int(binary.BigEndian.Uint64(v.Data))
	}
	if i <= 0 {
		return 0
	}
	return uint(i)
}

func (u uintConverter) FromInt8(v Variant) uint {
	var i int8
	if len(v.Data) == 1 {
		i = int8(v.Data[0])
	}
	if i <= 0 {
		return 0
	}
	return uint(i)
}

func (u uintConverter) FromInt16(v Variant) uint {
	var i int16
	if len(v.Data) == 2 {
		i = int16(binary.BigEndian.Uint16(v.Data))
	}
	if i <= 0 {
		return 0
	}
	return uint(i)
}

func (u uintConverter) FromInt32(v Variant) uint {
	i := int32(binary.BigEndian.Uint32(v.Data))
	if i <= 0 {
		return 0
	}
	return uint(i)
}

func (u uintConverter) FromInt64(v Variant) uint {
	i := int(binary.BigEndian.Uint64(v.Data))
	if i <= 0 {
		return 0
	}
	return uint(i)
}

func (u uintConverter) FromUint(v Variant) uint {
	var i uint
	if intSize == 32 {
		i = uint(binary.BigEndian.Uint32(v.Data))
	} else if intSize == 64 {
		i = uint(binary.BigEndian.Uint64(v.Data))
	}
	return i
}

func (u uintConverter) FromUint8(v Variant) uint {
	return uint(v.Data[0])
}

func (u uintConverter) FromUint16(v Variant) uint {
	return uint(binary.BigEndian.Uint16(v.Data))
}

func (u uintConverter) FromUint32(v Variant) uint {
	return uint(binary.BigEndian.Uint32(v.Data))
}

func (u uintConverter) FromUint64(v Variant) uint {
	return uint(binary.BigEndian.Uint64(v.Data))
}

func (u uintConverter) FromFloat32(v Variant) uint {
	var f float32
	if len(v.Data) == 4 {
		f = math.Float32frombits(binary.BigEndian.Uint32(v.Data))
	}
	return uint(f)
}

func (u uintConverter) FromFloat64(v Variant) uint {
	var f float64
	if len(v.Data) == 8 {
		f = math.Float64frombits(binary.BigEndian.Uint64(v.Data))
	}
	return uint(f)
}

func (u uintConverter) FromTime(v Variant) uint {
	var t time.Time
	err := t.UnmarshalBinary(v.Data)
	if err != nil {
		return 0
	}
	return uint(t.UnixNano())
}

func newUintConverter() *uintConverter {
	c := &uintConverter{}
	c.m = map[Kind]func(v Variant) uint{
		String:  c.FromString,
		Bool:    c.FromBool,
		Int:     c.FromInt,
		Int8:    c.FromInt8,
		Int16:   c.FromInt16,
		Int32:   c.FromInt32,
		Int64:   c.FromInt64,
		Uint:    c.FromUint,
		Uint8:   c.FromUint8,
		Uint16:  c.FromUint16,
		Uint32:  c.FromUint32,
		Uint64:  c.FromUint64,
		Float32: c.FromFloat32,
		Float64: c.FromFloat64,
		Time:    c.FromTime,
	}
	return c
}

var _ IConvertStrategy[float32] = (*float32Converter)(nil)

type float32Converter struct {
	Converter[float32]
}

func (c float32Converter) FromString(v Variant) float32 {
	s := *(*string)(unsafe.Pointer(&v.Data))
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f64)
}

func (c float32Converter) FromBool(v Variant) float32 {
	if len(v.Data) == 0 || v.Data[0] == 0x00 {
		return 0
	}
	return 1
}

func (c float32Converter) FromInt(v Variant) float32 {
	var i int
	if len(v.Data) == 4 && intSize == 32 {
		i = int(binary.BigEndian.Uint32(v.Data))
	} else if len(v.Data) == 8 && intSize == 64 {
		i = int(binary.BigEndian.Uint64(v.Data))
	}
	return float32(i)
}

func (c float32Converter) FromInt8(v Variant) float32 {
	var i int8
	if len(v.Data) == 1 {
		i = int8(v.Data[0])
	}
	return float32(i)
}

func (c float32Converter) FromInt16(v Variant) float32 {
	var i int16
	if len(v.Data) == 2 {
		i = int16(binary.BigEndian.Uint16(v.Data))
	}
	return float32(i)
}

func (c float32Converter) FromInt32(v Variant) float32 {
	return float32(int32(binary.BigEndian.Uint32(v.Data)))
}

func (c float32Converter) FromInt64(v Variant) float32 {
	return float32(int64(binary.BigEndian.Uint64(v.Data)))
}

func (c float32Converter) FromUint(v Variant) float32 {
	var i uint
	if intSize == 32 {
		i = uint(binary.BigEndian.Uint32(v.Data))
	} else if intSize == 64 {
		i = uint(binary.BigEndian.Uint64(v.Data))
	}
	return float32(i)
}

func (c float32Converter) FromUint8(v Variant) float32 {
	return float32(v.Data[0])
}

func (c float32Converter) FromUint16(v Variant) float32 {
	return float32(binary.BigEndian.Uint16(v.Data))
}

func (c float32Converter) FromUint32(v Variant) float32 {
	return float32(binary.BigEndian.Uint32(v.Data))
}

func (c float32Converter) FromUint64(v Variant) float32 {
	return float32(binary.BigEndian.Uint64(v.Data))
}

func (c float32Converter) FromFloat32(v Variant) float32 {
	var f float32
	if len(v.Data) == 4 {
		f = math.Float32frombits(binary.BigEndian.Uint32(v.Data))
	}
	return f
}

func (c float32Converter) FromFloat64(v Variant) float32 {
	var f float64
	if len(v.Data) == 8 {
		f = math.Float64frombits(binary.BigEndian.Uint64(v.Data))
	}
	return float32(f)
}

func (c float32Converter) FromTime(v Variant) float32 {
	var t time.Time
	err := t.UnmarshalBinary(v.Data)
	if err != nil {
		return 0
	}
	return float32(t.UnixNano())
}

func newFloat32Converter() *float32Converter {
	c := &float32Converter{}
	c.m = map[Kind]func(v Variant) float32{
		String:  c.FromString,
		Bool:    c.FromBool,
		Int:     c.FromInt,
		Int8:    c.FromInt8,
		Int16:   c.FromInt16,
		Int32:   c.FromInt32,
		Int64:   c.FromInt64,
		Uint:    c.FromUint,
		Uint8:   c.FromUint8,
		Uint16:  c.FromUint16,
		Uint32:  c.FromUint32,
		Uint64:  c.FromUint64,
		Float32: c.FromFloat32,
		Float64: c.FromFloat64,
		Time:    c.FromTime,
	}
	return c
}

var _ IConvertStrategy[float64] = (*float64Converter)(nil)

type float64Converter struct {
	Converter[float64]
}

func (c float64Converter) FromString(v Variant) float64 {
	s := *(*string)(unsafe.Pointer(&v.Data))
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f64
}

func (c float64Converter) FromBool(v Variant) float64 {
	if len(v.Data) == 0 || v.Data[0] == 0x00 {
		return 0
	}
	return 1
}

func (c float64Converter) FromInt(v Variant) float64 {
	var i int
	if len(v.Data) == 4 && intSize == 32 {
		i = int(binary.BigEndian.Uint32(v.Data))
	} else if len(v.Data) == 8 && intSize == 64 {
		i = int(binary.BigEndian.Uint64(v.Data))
	}
	return float64(i)
}

func (c float64Converter) FromInt8(v Variant) float64 {
	var i int8
	if len(v.Data) == 1 {
		i = int8(v.Data[0])
	}
	return float64(i)
}

func (c float64Converter) FromInt16(v Variant) float64 {
	var i int16
	if len(v.Data) == 2 {
		i = int16(binary.BigEndian.Uint16(v.Data))
	}
	return float64(i)
}

func (c float64Converter) FromInt32(v Variant) float64 {
	return float64(int32(binary.BigEndian.Uint32(v.Data)))
}

func (c float64Converter) FromInt64(v Variant) float64 {
	return float64(int64(binary.BigEndian.Uint64(v.Data)))
}

func (c float64Converter) FromUint(v Variant) float64 {
	var i uint
	if intSize == 32 {
		i = uint(binary.BigEndian.Uint32(v.Data))
	} else if intSize == 64 {
		i = uint(binary.BigEndian.Uint64(v.Data))
	}
	return float64(i)
}

func (c float64Converter) FromUint8(v Variant) float64 {
	return float64(v.Data[0])
}

func (c float64Converter) FromUint16(v Variant) float64 {
	return float64(binary.BigEndian.Uint16(v.Data))
}

func (c float64Converter) FromUint32(v Variant) float64 {
	return float64(binary.BigEndian.Uint32(v.Data))
}

func (c float64Converter) FromUint64(v Variant) float64 {
	return float64(binary.BigEndian.Uint64(v.Data))
}

func (c float64Converter) FromFloat32(v Variant) float64 {
	var f float32
	if len(v.Data) == 4 {
		f = math.Float32frombits(binary.BigEndian.Uint32(v.Data))
	}
	str := strconv.FormatFloat(float64(f), 'g', -1, 32)
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return f64
}

func (c float64Converter) FromFloat64(v Variant) float64 {
	var f float64
	if len(v.Data) == 8 {
		f = math.Float64frombits(binary.BigEndian.Uint64(v.Data))
	}
	return f
}

func (c float64Converter) FromTime(v Variant) float64 {
	var t time.Time
	err := t.UnmarshalBinary(v.Data)
	if err != nil {
		return 0
	}
	return float64(t.UnixNano())
}

func newFloat64Converter() *float64Converter {
	c := &float64Converter{}
	c.m = map[Kind]func(v Variant) float64{
		String:  c.FromString,
		Bool:    c.FromBool,
		Int:     c.FromInt,
		Int8:    c.FromInt8,
		Int16:   c.FromInt16,
		Int32:   c.FromInt32,
		Int64:   c.FromInt64,
		Uint:    c.FromUint,
		Uint8:   c.FromUint8,
		Uint16:  c.FromUint16,
		Uint32:  c.FromUint32,
		Uint64:  c.FromUint64,
		Float32: c.FromFloat32,
		Float64: c.FromFloat64,
		Time:    c.FromTime,
	}
	return c
}

var _ IConvertStrategy[time.Time] = (*timeConverter)(nil)

type timeConverter struct {
	Converter[time.Time]
}

// FromBool implements IConvertStrategy.
func (t *timeConverter) FromBool(v Variant) time.Time {
	if len(v.Data) == 0 || v.Data[0] == 0x00 {
		return time.Time{}
	}
	return time.Now()
}

// FromFloat32 implements IConvertStrategy.
func (t *timeConverter) FromFloat32(v Variant) time.Time {
	var f float32
	if len(v.Data) == 4 {
		f = math.Float32frombits(binary.BigEndian.Uint32(v.Data))
	}
	return time.Unix(0, int64(f))
}

// FromFloat64 implements IConvertStrategy.
func (t *timeConverter) FromFloat64(v Variant) time.Time {
	var f float64
	if len(v.Data) == 8 {
		f = math.Float64frombits(binary.BigEndian.Uint64(v.Data))
	}
	return time.Unix(0, int64(f))
}

// FromInt implements IConvertStrategy.
func (t *timeConverter) FromInt(v Variant) time.Time {
	var i int
	if len(v.Data) == 4 && intSize == 32 {
		i = int(binary.BigEndian.Uint32(v.Data))
	} else if len(v.Data) == 8 && intSize == 64 {
		i = int(binary.BigEndian.Uint64(v.Data))
	}
	return time.Unix(0, int64(i))
}

// FromInt16 implements IConvertStrategy.
func (t *timeConverter) FromInt16(v Variant) time.Time {
	var i int16
	if len(v.Data) == 2 {
		i = int16(binary.BigEndian.Uint16(v.Data))
	}
	return time.Unix(0, int64(i))
}

// FromInt32 implements IConvertStrategy.
func (t *timeConverter) FromInt32(v Variant) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint32(v.Data)))
}

// FromInt64 implements IConvertStrategy.
func (t *timeConverter) FromInt64(v Variant) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint64(v.Data)))
}

// FromInt8 implements IConvertStrategy.
func (t *timeConverter) FromInt8(v Variant) time.Time {
	var i int8
	if len(v.Data) == 1 {
		i = int8(v.Data[0])
	}
	return time.Unix(0, int64(i))
}

// FromString implements IConvertStrategy.
func (t *timeConverter) FromString(v Variant) time.Time {
	s := *(*string)(unsafe.Pointer(&v.Data))
	tt, err := time.Parse(v.layout, s)
	if err != nil {
	}
	return tt
}

// FromTime implements IConvertStrategy.
func (t *timeConverter) FromTime(v Variant) time.Time {
	var tt time.Time
	err := tt.UnmarshalBinary(v.Data)
	if err != nil {
		return time.Time{}
	}
	return tt
}

// FromUint implements IConvertStrategy.
func (t *timeConverter) FromUint(v Variant) time.Time {
	var i uint
	if intSize == 32 {
		i = uint(binary.BigEndian.Uint32(v.Data))
	} else if intSize == 64 {
		i = uint(binary.BigEndian.Uint64(v.Data))
	}
	return time.Unix(0, int64(i))
}

// FromUint16 implements IConvertStrategy.
func (t *timeConverter) FromUint16(v Variant) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint16(v.Data)))
}

// FromUint32 implements IConvertStrategy.
func (t *timeConverter) FromUint32(v Variant) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint32(v.Data)))
}

// FromUint64 implements IConvertStrategy.
func (t *timeConverter) FromUint64(v Variant) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint64(v.Data)))
}

// FromUint8 implements IConvertStrategy.
func (t *timeConverter) FromUint8(v Variant) time.Time {
	return time.Unix(0, int64(v.Data[0]))
}

func newTimeConverter() *timeConverter {
	c := &timeConverter{}
	c.m = map[Kind]func(v Variant) time.Time{
		String:  c.FromString,
		Bool:    c.FromBool,
		Int:     c.FromInt,
		Int8:    c.FromInt8,
		Int16:   c.FromInt16,
		Int32:   c.FromInt32,
		Int64:   c.FromInt64,
		Uint:    c.FromUint,
		Uint8:   c.FromUint8,
		Uint16:  c.FromUint16,
		Uint32:  c.FromUint32,
		Uint64:  c.FromUint64,
		Float32: c.FromFloat32,
		Float64: c.FromFloat64,
		Time:    c.FromTime,
	}
	return c
}
