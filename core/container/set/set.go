package set

// note: visit https://github.dev/deckarep/golang-set for more details.

type Set[T comparable] interface {
	// Add adds an element to the set.
	Add(v T) bool

	// Append multiple elements to the set.
	Append(vals ...T) int

	// Contains reports whether a single value is present in set.
	Contains(v T) bool

	// ContainsAll reports whether values are all present in set.
	ContainsAll(vals ...T) bool

	// ContainsAny reports whether at least one of the values is
	// present in set.
	ContainsAny(vals ...T) bool

	// Size returns the number of elements in the set.
	Size() int

	// IsEmpty reports whether the set is empty.
	IsEmpty() bool

	// Pop an arbitrary element in the set.
	Pop() (v T, ok bool)

	// Remove the given element from the set.
	Remove(v T)

	// Clear all the elements.
	Clear()

	// Clone returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Set[T]

	// Difference returns the difference between this set
	// and other.
	Difference(other Set[T]) Set[T]

	// SymmetricDifference returns a new set with all elements which are
	// in either this set or the other set but not in both.
	SymmetricDifference(other Set[T]) Set[T]

	// Intersect returns a new set containing only the elements
	// that exist only in both sets.
	Intersect(other Set[T]) Set[T]

	// Union returns a new set with all elements in both sets.
	Union(other Set[T]) Set[T]

	// ToSlice returns the members of the set as a slice.
	ToSlice() []T

	// String provides a convenient string representation
	// of the current state of the set.
	String() string

	// MarshalJSON will marshal the set into a JSON-based representation.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON will unmarshal a JSON-based byte slice into a full Set datastructure.
	// For this to work, set subtypes must implemented the Marshal/Unmarshal interface.
	UnmarshalJSON(b []byte) error
}

// NewUnsafeSet creates a thread unsafe set.
func NewUnsafeSet[T comparable](vals ...T) Set[T] {
	s := make(set[T], len(vals))
	for _, item := range vals {
		s.add(item)
	}
	return s
}

// NewSet creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewSet[T comparable](vals ...T) Set[T] {
	s := &safeSet[T]{
		set: make(set[T], len(vals)),
	}
	for _, item := range vals {
		s.Add(item)
	}
	return s
}
