package variant

import (
	"fmt"
	"testing"
)

func BenchmarkVariant(b *testing.B) {
	v := New(int64(-128))
	for i := 0; i < b.N; i++ {
		fmt.Println(v.ToString())
	}
	b.ReportAllocs()
}
