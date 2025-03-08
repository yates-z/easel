package variant

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"
)

const (
	intSize = 32 << (^uint(0) >> 63)
	maxUint = 1<<intSize - 1
	maxInt  = 1<<(intSize-1) - 1
)

type Variant struct {
	Type   Kind
	Data   []byte
	layout string
}

var Nil = Variant{Type: Invalid, Data: make([]byte, 0, 8)}

func (v *Variant) SetLayout(layout string) *Variant {
	v.layout = layout
	return v
}

func (v Variant) String() string {
	return fmt.Sprintf("Variant(%v, %v)", v.Type, v.Data)
}

func (v Variant) ToString() string {

	if fn := Strategies.string.Get(v.Type); fn != nil {
		return fn(v)
	}
	return ""
}

func (v Variant) ToInt() int {
	if fn := Strategies.int.Get(v.Type); fn != nil {
		return fn(v)
	}
	return 0
}

func (v Variant) ToBytes() []byte {
	return v.Data
}

func (v Variant) ToUint() uint {
	if fn := Strategies.uint.Get(v.Type); fn != nil {
		return fn(v)
	}
	return 0
}

func (v Variant) ToBool() bool {
	for _, ch := range v.Data {
		if ch != 0 {
			return true
		}
	}
	return false
}

func (v Variant) ToFloat32() float32 {
	if fn := Strategies.float32.Get(v.Type); fn != nil {
		return fn(v)
	}
	return 0
}

func (v Variant) ToFloat64() float64 {
	if fn := Strategies.float64.Get(v.Type); fn != nil {
		return fn(v)
	}
	return 0
}

func (v Variant) ToTime() time.Time {
	if fn := Strategies.time.Get(v.Type); fn != nil {
		return fn(v)
	}
	return time.Time{}
}

func (v Variant) Equal(other any) bool {
	if v, ok := other.(Variant); ok {
		return reflect.DeepEqual(v, other)
	}
	variant := New(other)
	variant.layout = v.layout
	return reflect.DeepEqual(v, variant)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Variant) UnmarshalJSON(data []byte) error {
	v.Type = String
	if data[0] == '"' && data[len(data)-1] == '"' {
		v.Data = data[1 : len(data)-1]
	} else if data[0] == 't' {
		v.Type = Bool
		v.Data = append(v.Data, 0x01)
	} else if data[0] == 'f' {
		v.Type = Bool
		v.Data = append(v.Data, 0x00)
	} else if data[0] == 'n' {
		v.Type = Invalid
	} else {
		v.Data = data
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v Variant) MarshalJSON() ([]byte, error) {
	return v.Data, nil
}

func New(v any) Variant {
	variant := Variant{
		Type:   Invalid,
		Data:   make([]byte, 0, 8),
		layout: "2006-01-02T15:04:05Z",
	}

	switch v := v.(type) {
	case string:
		variant.Type = String
		variant.Data = append(variant.Data, v...)
	case bool:
		variant.Type = Bool
		if v {
			variant.Data = append(variant.Data, 0x01)
		} else {
			variant.Data = append(variant.Data, 0x00)
		}
	case int:
		variant.Type = Int
		if intSize == 32 {
			variant.Data = append(variant.Data, make([]byte, 4)...)
			binary.BigEndian.PutUint32(variant.Data, uint32(v))
		} else if intSize == 64 {
			variant.Data = append(variant.Data, make([]byte, 8)...)
			binary.BigEndian.PutUint64(variant.Data, uint64(v))
		}
	case int8:
		variant.Type = Int8
		variant.Data = append(variant.Data, byte(v))
	case int16:
		variant.Type = Int16
		variant.Data = append(variant.Data, make([]byte, 2)...)
		binary.BigEndian.PutUint16(variant.Data, uint16(v))
	case int32:
		variant.Type = Int32
		variant.Data = append(variant.Data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(variant.Data, uint32(v))
	case int64:
		variant.Type = Int64
		variant.Data = append(variant.Data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(variant.Data, uint64(v))
	case uint:
		variant.Type = Uint
		if intSize == 32 {
			variant.Data = append(variant.Data, make([]byte, 4)...)
			binary.BigEndian.PutUint32(variant.Data, uint32(v))
		} else if intSize == 64 {
			variant.Data = append(variant.Data, make([]byte, 8)...)
			binary.BigEndian.PutUint64(variant.Data, uint64(v))
		}
	case uint8:
		variant.Type = Uint8
		variant.Data = append(variant.Data, v)
	case uint16:
		variant.Type = Uint16
		variant.Data = append(variant.Data, make([]byte, 2)...)
		binary.BigEndian.PutUint16(variant.Data, v)
	case uint32:
		variant.Type = Uint32
		variant.Data = append(variant.Data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(variant.Data, v)
	case uint64:
		variant.Type = Uint64
		variant.Data = append(variant.Data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(variant.Data, v)
	case float32:
		variant.Type = Float32
		variant.Data = append(variant.Data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(variant.Data, math.Float32bits(v))
	case float64:
		variant.Type = Float64
		variant.Data = append(variant.Data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(variant.Data, math.Float64bits(v))
	case time.Time:
		data, err := v.MarshalBinary()
		if err == nil {
			variant.Type = Time
			variant.Data = data
		}
	}
	return variant
}
