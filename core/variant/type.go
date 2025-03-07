package variant

import "strconv"

type Kind uint8

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Interface
	String
	Time
)

func (k Kind) String() string {
	if uint(k) < uint(len(kindNames)) {
		return kindNames[uint(k)]
	}
	return "kind" + strconv.Itoa(int(k))
}

var kindNames = []string{
	Invalid:    "invalid",
	Bool:       "bool",
	Int:        "int",
	Int8:       "int8",
	Int16:      "int16",
	Int32:      "int32",
	Int64:      "int64",
	Uint:       "uint",
	Uint8:      "uint8",
	Uint16:     "uint16",
	Uint32:     "uint32",
	Uint64:     "uint64",
	Uintptr:    "uintptr",
	Float32:    "float32",
	Float64:    "float64",
	Complex64:  "complex64",
	Complex128: "complex128",
	Interface:  "interface",
	String:     "string",
	Time:       "time.Time",
}
