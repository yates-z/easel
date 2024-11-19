package tuple

import (
	"github.com/yates-z/easel/core/variant"
)

type tuple struct {
	elements []variant.Variant
}

func NewTuple(elems ...any) tuple {
	t := tuple{elements: make([]variant.Variant, 0, len(elems))}
	for _, elem := range elems {
		t.elements = append(t.elements, variant.New(elem))
	}
	return t
}

func (t tuple) Get(index int) (variant.Variant, bool) {
	if index < 0 || index >= len(t.elements) {
		var value variant.Variant
		return value, false
	}
	return t.elements[index], true
}

func (t tuple) Count(val any) int {
	count := 0
	for _, elem := range t.elements {
		if elem.Equal(val) {
			count++
		}
	}
	return count
}

func (t tuple) Len() int {
	return len(t.elements)
}

func (t tuple) ToSlice() []variant.Variant {
	copied := make([]variant.Variant, len(t.elements))
	copy(copied, t.elements)

	return copied
}
