package set

import "sync"

var _ Set[string] = (*safeSet[string])(nil)

type safeSet[T comparable] struct {
	sync.RWMutex
	set set[T]
}

func (s *safeSet[T]) Add(v T) bool {
	s.Lock()
	ret := s.set.Add(v)
	s.Unlock()
	return ret
}

func (s *safeSet[T]) Append(vals ...T) int {
	s.Lock()
	ret := s.set.Append(vals...)
	s.Unlock()
	return ret
}

func (s *safeSet[T]) Contains(v T) bool {
	s.RLock()
	ret := s.set.Contains(v)
	s.RUnlock()

	return ret
}

func (s *safeSet[T]) ContainsAll(vals ...T) bool {
	s.RLock()
	ret := s.set.ContainsAll(vals...)
	s.RUnlock()

	return ret
}

func (s *safeSet[T]) ContainsAny(vals ...T) bool {
	s.RLock()
	ret := s.set.ContainsAny(vals...)
	s.RUnlock()

	return ret
}

func (s *safeSet[T]) Size() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.set)
}

func (s *safeSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *safeSet[T]) Pop() (T, bool) {
	s.Lock()
	defer s.Unlock()
	return s.set.Pop()
}

func (s *safeSet[T]) Clear() {
	s.Lock()
	s.set.Clear()
	s.Unlock()
}

func (s *safeSet[T]) Remove(v T) {
	s.Lock()
	delete(s.set, v)
	s.Unlock()
}

func (s *safeSet[T]) Clone() Set[T] {
	s.RLock()

	unsafeClone := s.set.Clone().(set[T])
	ret := &safeSet[T]{set: unsafeClone}
	s.RUnlock()
	return ret
}

func (s *safeSet[T]) Union(other Set[T]) Set[T] {
	o := other.(*safeSet[T])

	s.RLock()
	o.RLock()

	unsafeUnion := s.set.Union(o.set).(set[T])
	ret := &safeSet[T]{set: unsafeUnion}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet[T]) Intersect(other Set[T]) Set[T] {
	o := other.(*safeSet[T])

	s.RLock()
	o.RLock()

	unsafeIntersection := s.set.Intersect(o.set).(set[T])
	ret := &safeSet[T]{set: unsafeIntersection}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet[T]) Difference(other Set[T]) Set[T] {
	o := other.(*safeSet[T])

	s.RLock()
	o.RLock()

	unsafeDifference := s.set.Difference(o.set).(set[T])
	ret := &safeSet[T]{set: unsafeDifference}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet[T]) SymmetricDifference(other Set[T]) Set[T] {
	o := other.(*safeSet[T])

	s.RLock()
	o.RLock()

	unsafeDifference := s.set.SymmetricDifference(o.set).(set[T])
	ret := &safeSet[T]{set: unsafeDifference}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet[T]) ToSlice() []T {
	keys := make([]T, 0, s.Size())
	s.RLock()
	for elem := range s.set {
		keys = append(keys, elem)
	}
	s.RUnlock()
	return keys
}

func (s *safeSet[T]) String() string {
	s.RLock()
	ret := s.set.String()
	s.RUnlock()
	return ret
}

func (s *safeSet[T]) MarshalJSON() ([]byte, error) {
	s.RLock()
	b, err := s.set.MarshalJSON()
	s.RUnlock()

	return b, err
}

func (s *safeSet[T]) UnmarshalJSON(p []byte) error {
	s.RLock()
	err := s.set.UnmarshalJSON(p)
	s.RUnlock()

	return err
}
