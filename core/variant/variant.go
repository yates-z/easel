package variant

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

const intSize = 32 << (^uint(0) >> 63)

type Variant struct {
	Type Kind
	Data []byte
}

func (v Variant) String() string {
	return fmt.Sprintf("Variant(%v, %v)", v.Type, v.Data)
}

func (v Variant) ToString() string {
	switch v.Type {
	case String:
		return *(*string)(unsafe.Pointer(&v.Data))
	case Bool:
		if v.Data[0] == 0x00 {
			return "false"
		} else {
			return "true"
		}
	case Int:
		var i int
		if intSize == 32 {
			i = int(binary.BigEndian.Uint32(v.Data))
		} else if intSize == 64 {
			i = int(binary.BigEndian.Uint64(v.Data))
		}
		return strconv.Itoa(i)
	case Int8:
		i := v.Data[0]
		s := strconv.AppendInt(make([]byte, 0), int64(int8(i)), 10)
		return *(*string)(unsafe.Pointer(&s))
	case Int16:
		i := int16(binary.BigEndian.Uint16(v.Data))
		return strconv.Itoa(int(i))
	case Int32:
		i := int32(binary.BigEndian.Uint32(v.Data))
		return strconv.Itoa(int(i))
	case Int64:
		i := int64(binary.BigEndian.Uint64(v.Data))
		return strconv.Itoa(int(i))
	case Uint:
		var i uint64
		if intSize == 32 {
			i = uint64(binary.BigEndian.Uint32(v.Data))
		} else if intSize == 64 {
			i = binary.BigEndian.Uint64(v.Data)
		}
		return strconv.FormatUint(i, 10)
	case Uint8:
		i := v.Data[0]
		s := strconv.AppendUint(make([]byte, 0), uint64(i), 10)
		return *(*string)(unsafe.Pointer(&s))
	case Uint16:
		i := uint64(binary.BigEndian.Uint16(v.Data))
		return strconv.FormatUint(i, 10)
	case Uint32:
		i := uint64(binary.BigEndian.Uint32(v.Data))
		return strconv.FormatUint(i, 10)
	case Uint64:
		i := binary.BigEndian.Uint64(v.Data)
		return strconv.FormatUint(i, 10)
	case Float32:
		f := math.Float32frombits(binary.BigEndian.Uint32(v.Data))
		return strconv.FormatFloat(float64(f), 'g', -1, 32)
	case Float64:
		f := math.Float64frombits(binary.BigEndian.Uint64(v.Data))
		return strconv.FormatFloat(f, 'g', -1, 64)
	default:
		return ""
	}
}

func (v Variant) ToBytes() string {
	return ""
}

func New(v any) Variant {
	variant := Variant{
		Type: Invalid,
		Data: make([]byte, 0, 8),
	}

	switch v.(type) {
	case string:
		variant.Type = String
		variant.Data = append(variant.Data, v.(string)...)
	case bool:
		variant.Type = Bool
		if v.(bool) {
			variant.Data = append(variant.Data, 0x01)
		} else {
			variant.Data = append(variant.Data, 0x00)
		}
	case int:
		variant.Type = Int
		if intSize == 32 {
			variant.Data = append(variant.Data, make([]byte, 4)...)
			binary.BigEndian.PutUint32(variant.Data, uint32(v.(int)))
		} else if intSize == 64 {
			variant.Data = append(variant.Data, make([]byte, 8)...)
			binary.BigEndian.PutUint64(variant.Data, uint64(v.(int)))
		}
	case int8:
		variant.Type = Int8
		variant.Data = append(variant.Data, byte(v.(int8)))
	case int16:
		variant.Type = Int16
		variant.Data = append(variant.Data, make([]byte, 2)...)
		binary.BigEndian.PutUint16(variant.Data, uint16(v.(int16)))
	case int32:
		variant.Type = Int32
		variant.Data = append(variant.Data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(variant.Data, uint32(v.(int32)))
	case int64:
		variant.Type = Int64
		variant.Data = append(variant.Data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(variant.Data, uint64(v.(int64)))
	case uint:
		variant.Type = Uint
		if intSize == 32 {
			variant.Data = append(variant.Data, make([]byte, 4)...)
			binary.BigEndian.PutUint32(variant.Data, uint32(v.(uint)))
		} else if intSize == 64 {
			variant.Data = append(variant.Data, make([]byte, 8)...)
			binary.BigEndian.PutUint64(variant.Data, uint64(v.(uint)))
		}
	case uint8:
		variant.Type = Uint8
		variant.Data = append(variant.Data, v.(uint8))
	case uint16:
		variant.Type = Uint16
		variant.Data = append(variant.Data, make([]byte, 2)...)
		binary.BigEndian.PutUint16(variant.Data, v.(uint16))
	case uint32:
		variant.Type = Uint32
		variant.Data = append(variant.Data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(variant.Data, v.(uint32))
	case uint64:
		variant.Type = Uint64
		variant.Data = append(variant.Data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(variant.Data, v.(uint64))
	case float32:
		variant.Type = Float32
		variant.Data = append(variant.Data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(variant.Data, math.Float32bits(v.(float32)))
	case float64:
		variant.Type = Float64
		variant.Data = append(variant.Data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(variant.Data, math.Float64bits(v.(float64)))
	}
	return variant
}
