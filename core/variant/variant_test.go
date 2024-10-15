package variant

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func assert(condition bool, message string) {
	if !condition {
		log.Fatal(message)
	}
}

func BenchmarkVariant(b *testing.B) {
	v := New("-128")
	for i := 0; i < b.N; i++ {
		fmt.Println(v.ToInt())
	}
	b.ReportAllocs()
}

func TestVariant_ToInt(t *testing.T) {
	// bool
	v := New(true)
	assert(v.ToInt() == 1, "")
	v = New(false)
	assert(v.ToInt() == 0, "")
	v = New(-9223372036854775808)
	assert(v.ToInt() == -9223372036854775808, "")

	v = New(uint(9223372036854775810))
	assert(v.ToInt() == 0, "")

	v = New(int8(-100))
	assert(v.ToInt() == -100, "")
	v = New(int8(100))
	assert(v.ToInt() == 100, "")
	v = New(int16(-100))
	assert(v.ToInt() == -100, "")
	v = New(int32(-100))
	assert(v.ToInt() == -100, "")
	v = New(int64(-100))
	assert(v.ToInt() == -100, "")
	v = New(uint8(255))
	assert(v.ToInt() == 255, "")
	v = New(uint16(65535))
	assert(v.ToInt() == 65535, "")

	v = New(-65535)
	assert(v.ToInt() == -65535, "")
	v = New(100.86)
	assert(v.ToInt() == 100, "")
	v = New(float32(100.86))
	assert(v.ToInt() == 100, "")
	v = New(float32(-100.86))
	assert(v.ToInt() == -100, "")
}

func TestVariant_ToString(t *testing.T) {
	v := New("hello world")
	assert(v.ToString() == "hello world", "")
	v = New(int8(-100))
	assert(v.ToString() == "-100", "")
	v = New(int16(-100))
	assert(v.ToString() == "-100", "")
	v = New(int32(-100))
	assert(v.ToString() == "-100", "")
	v = New(int64(-100))
	assert(v.ToString() == "-100", "")
	v = New(-100)
	assert(v.ToString() == "-100", "")
	s := -255
	v = New(uint8(s))
	assert(v.ToString() == "1", "")
	v = New(float32(100.86))
	assert(v.ToString() == "100.86", "")
	v = New(100.86)
	assert(v.ToString() == "100.86", "")
	v = New(-100.86)
	assert(v.ToString() == "-100.86", "")
	v = New(true)
	assert(v.ToString() == "true", "")
	v = New(false)
	assert(v.ToString() == "false", "")
	v = New(map[string]int{})
	assert(v.ToString() == "", "")
}

func TestVariant_ToUint(t *testing.T) {
	// bool
	v := New(true)
	assert(v.ToUint() == 1, "")
	v = New(false)
	assert(v.ToUint() == 0, "")
	v = New(-9223372036854775808)
	assert(v.ToUint() == 0, "")

	v = New(uint(9223372036854775810))
	assert(v.ToUint() == 9223372036854775810, "")

	v = New(int8(-100))
	assert(v.ToUint() == 0, "")
	v = New(int8(100))
	assert(v.ToUint() == 100, "")
	v = New(int16(-100))
	assert(v.ToUint() == 0, "")
	v = New(int32(-100))
	assert(v.ToUint() == 0, "")
	v = New(int64(-100))
	assert(v.ToUint() == 0, "")
	v = New(uint8(255))
	assert(v.ToUint() == 255, "")
	v = New(uint16(65535))
	assert(v.ToUint() == 65535, "")

	v = New(-65535)
	assert(v.ToUint() == 0, "")
	v = New(100.86)
	assert(v.ToUint() == 100, "")
	v = New(float32(100.86))
	assert(v.ToUint() == 100, "")
	v = New(float32(-100.86))
	assert(v.ToUint() == 0, "")
}

func TestVariant_ToBool(t *testing.T) {
	// bool
	v := New(true)
	assert(v.ToBool() == true, "")
	v = New(false)
	assert(v.ToBool() == false, "")
	v = New(uint(0))
	assert(v.ToBool() == false, "")
	v = New(uint(86))
	assert(v.ToBool() == true, "")
	v = New(uint8(0))
	assert(v.ToBool() == false, "")
	v = New(uint8(86))
	assert(v.ToBool() == true, "")
	v = New(uint16(0))
	assert(v.ToBool() == false, "")
	v = New(uint16(10086))
	assert(v.ToBool() == true, "")
	v = New(uint32(0))
	assert(v.ToBool() == false, "")
	v = New(uint32(10086))
	assert(v.ToBool() == true, "")
	v = New(uint64(0))
	assert(v.ToBool() == false, "")
	v = New(uint64(10086))
	assert(v.ToBool() == true, "")
	v = New("")
	assert(v.ToBool() == false, "")
	v = New("hello world")
	assert(v.ToBool() == true, "")

	v = New(0)
	assert(v.ToBool() == false, "")
	v = New(1)
	assert(v.ToBool() == true, "")
	v = New(int8(0))
	assert(v.ToBool() == false, "")
	v = New(int8(-1))
	assert(v.ToBool() == true, "")
	v = New(int16(0))
	assert(v.ToBool() == false, "")
	v = New(int16(-1))
	assert(v.ToBool() == true, "")
	v = New(int32(0))
	assert(v.ToBool() == false, "")
	v = New(int32(-1))
	assert(v.ToBool() == true, "")
	v = New(int64(0))
	assert(v.ToBool() == false, "")
	v = New(int64(-1))
	assert(v.ToBool() == true, "")
}

func TestVariant_ToFloat32(t *testing.T) {
	v := New(true)
	assert(v.ToFloat32() == 1, "")
	v = New(false)
	assert(v.ToFloat32() == 0, "")
	v = New(uint(0))
	assert(v.ToFloat32() == 0, "")
	v = New(uint8(86))
	assert(v.ToFloat32() == 86, "")
	v = New(10086)
	assert(v.ToFloat32() == 10086, "")
	v = New(int8(86))
	assert(v.ToFloat32() == 86, "")
	v = New(int16(0))
	assert(v.ToFloat32() == 0, "")
	v = New(int32(86))
	assert(v.ToFloat32() == 86, "")
	v = New(int64(86))
	assert(v.ToFloat32() == 86, "")
	v = New("100.86")
	assert(v.ToFloat32() == 100.86, "")
	v = New("-100.86")
	assert(v.ToFloat32() == -100.86, "")
	v = New("1abc")
	assert(v.ToFloat32() == 0, "")
	v = New(float32(-86.1))
	assert(v.ToFloat32() == -86.1, "")
	v = New(-86.1)
	assert(v.ToFloat32() == -86.1, "")
}

func TestVariant_ToFloat64(t *testing.T) {
	v := New(true)
	assert(v.ToFloat64() == 1, "1")
	v = New(false)
	assert(v.ToFloat64() == 0, "2")
	v = New(uint(0))
	assert(v.ToFloat64() == 0, "3")
	v = New(uint8(86))
	assert(v.ToFloat64() == 86, "4")
	v = New(10086)
	assert(v.ToFloat64() == 10086, "5")
	v = New(int8(86))
	assert(v.ToFloat64() == 86, "6")
	v = New(int16(0))
	assert(v.ToFloat64() == 0, "7")
	v = New(int32(86))
	assert(v.ToFloat64() == 86, "8")
	v = New(int64(86))
	assert(v.ToFloat64() == 86, "9")
	v = New("100.86")
	assert(v.ToFloat64() == 100.86, "10")
	v = New("-100.86")
	assert(v.ToFloat64() == -100.86, "11")
	v = New("1abc")
	assert(v.ToFloat64() == 0, "12")
	v = New(float32(-86.2))
	assert(v.ToFloat64() == -86.2, "13")
	v = New(-86.1)
	assert(v.ToFloat64() == -86.1, "14")
}

func TestVariant_Empty(t *testing.T) {
	v := Variant{Type: Bool}
	assert(v.ToString() == "false", "")
	assert(!v.ToBool(), "")
	assert(v.ToInt() == 0, "")
	assert(v.ToUint() == 0, "")
	assert(v.ToFloat32() == 0, "")
	assert(v.ToFloat64() == 0, "")

	v = Variant{Type: String}
	assert(v.ToString() == "", "1")
	assert(!v.ToBool(), "2")
	assert(v.ToInt() == 0, "3")
	assert(v.ToUint() == 0, "4")
	assert(v.ToFloat32() == 0, "5")
	assert(v.ToFloat64() == 0, "6")

	v = Variant{Type: Int}
	assert(v.ToString() == "0", "")
	assert(!v.ToBool(), "")
	assert(v.ToInt() == 0, "")
	assert(v.ToUint() == 0, "")
	assert(v.ToFloat32() == 0, "")
	assert(v.ToFloat64() == 0, "")

	v = Variant{Type: Float32}
	assert(v.ToString() == "0", "")
	assert(!v.ToBool(), "")
	assert(v.ToInt() == 0, "")
	assert(v.ToUint() == 0, "")
	assert(v.ToFloat32() == 0, "")
	assert(v.ToFloat64() == 0, "")
	v = Variant{Type: Float64}
	assert(v.ToString() == "0", "")
	assert(!v.ToBool(), "")
	assert(v.ToInt() == 0, "")
	assert(v.ToUint() == 0, "")
	assert(v.ToFloat32() == 0, "")
	assert(v.ToFloat64() == 0, "")
}

func Test_UnmarshalJson(t *testing.T) {
	data := []byte(`{"hello": "world!", "age": 36, "weight": 60.123, "is_male": null}`)
	m := map[string]Variant{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m["hello"].ToString())
	fmt.Println(m["age"].ToString())
	fmt.Println(m["weight"].ToString())
	fmt.Println(m["is_male"].ToBool())
	fmt.Println(m)
}
