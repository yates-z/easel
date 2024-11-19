package set

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ Set[string] = (set[string])(nil)

type set[T comparable] map[T]struct{}

func (s set[T]) add(v T) {
	s[v] = struct{}{}
}

func (s set[T]) Add(v T) bool {
	prevLen := len(s)
	s.add(v)
	return prevLen != len(s)
}

func (s set[T]) Append(vals ...T) int {
	prevLen := len(s)
	for _, val := range vals {
		s.add(val)
	}
	return len(s) - prevLen
}

func (s set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s set[T]) ContainsAll(vals ...T) bool {
	for _, val := range vals {
		if _, ok := s[val]; !ok {
			return false
		}
	}
	return true
}

func (s set[T]) ContainsAny(vals ...T) bool {
	for _, val := range vals {
		if _, ok := s[val]; ok {
			return true
		}
	}
	return false
}

func (s set[T]) Size() int {
	return len(s)
}

func (s set[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s set[T]) Pop() (v T, ok bool) {
	for item := range s {
		delete(s, item)
		return item, true
	}
	return v, false
}

func (s set[T]) Remove(v T) {
	delete(s, v)
}

func (s set[T]) Clear() {
	for key := range s {
		delete(s, key)
	}
}

func (s set[T]) Clone() Set[T] {
	clonedSet := make(set[T], len(s))
	for elem := range s {
		clonedSet.Add(elem)
	}
	return clonedSet
}

func (s set[T]) Difference(other Set[T]) Set[T] {
	o := other.(set[T])

	diff := make(set[T])
	for elem := range s {
		if !o.Contains(elem) {
			diff.add(elem)
		}
	}
	return diff
}

func (s set[T]) SymmetricDifference(other Set[T]) Set[T] {
	o := other.(set[T])

	sd := make(set[T])
	for elem := range s {
		if !o.Contains(elem) {
			sd.add(elem)
		}
	}
	for elem := range o {
		if !s.Contains(elem) {
			sd.add(elem)
		}
	}
	return sd
}

func (s set[T]) Intersect(other Set[T]) Set[T] {
	o := other.(set[T])

	intersection := make(set[T])
	// loop over smaller set
	if s.Size() < other.Size() {
		for elem := range s {
			if o.Contains(elem) {
				intersection.add(elem)
			}
		}
	} else {
		for elem := range o {
			if s.Contains(elem) {
				intersection.add(elem)
			}
		}
	}
	return intersection
}

func (s set[T]) Union(other Set[T]) Set[T] {
	o := other.(set[T])

	n := s.Size()
	if o.Size() > n {
		n = o.Size()
	}
	unionedSet := make(set[T], n)

	for elem := range s {
		unionedSet.add(elem)
	}
	for elem := range o {
		unionedSet.add(elem)
	}
	return unionedSet
}

func (s set[T]) ToSlice() []T {
	keys := make([]T, 0, s.Size())
	for elem := range s {
		keys = append(keys, elem)
	}

	return keys
}

func (s set[T]) String() string {
	items := make([]string, 0, len(s))

	for elem := range s {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (s set[T]) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, s.Size())

	for elem := range s {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON recreates a set from a JSON array, it only decodes
// primitive types. Numbers are decoded as json.Number.
func (s set[T]) UnmarshalJSON(b []byte) error {
	var i []T
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	s.Append(i...)

	return nil
}
