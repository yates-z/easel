package tuple

import (
	"fmt"
	"testing"
)

func Test_Tuple(t *testing.T) {

	tuple := NewTuple("hello", "wolrd")
	v, ok := tuple.Get(0)
	fmt.Println(v.ToString(), ok)
	v, ok = tuple.Get(1)
	fmt.Println(v.ToString(), ok)
	v, ok = tuple.Get(2)
	fmt.Println(v.ToString(), ok)

	tuple = NewTuple("1", "2", "hello", "2", "hello", "2")

	fmt.Println("Count of 2:", tuple.Count("2"))           // 输出: Count of 2: 3
	fmt.Println("Count of 'hello':", tuple.Count("hello")) // 输出: Count of 'hello': 2
	fmt.Println("Count of 5:", tuple.Count("5"))

	s := tuple.ToSlice()
	fmt.Println(s)

}
