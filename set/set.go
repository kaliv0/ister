package set

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

type Set[T comparable] struct {
	data map[T]struct{}
}

// Empty set. Zero-capacity backing map; first Add allocates.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{}
}

// Set from the given values. Duplicates are dropped. SetOf[T]() with no args is empty.
func SetOf[T comparable](vals ...T) *Set[T] {
	s := NewSet[T]()
	s.Add(vals...)
	return s
}

// Set from a slice. Copy `s` so later mutations of the original slice do not affect the set.
//
//	Duplicates are dropped. `SetFromSlice(nil)` is empty.
func SetFromSlice[T comparable](s []T) *Set[T] {
	return SetOf(s...)
}

// Number of elements.
func (s *Set[T]) Len() int {
	return len(s.data)
}

// Add one or more elements.
func (s *Set[T]) Add(vals ...T) {
	if len(vals) == 0 {
		return
	}

	if s.data == nil {
		s.data = make(map[T]struct{}, len(vals))
	}

	for _, v := range vals {
		s.data[v] = struct{}{}
	}
}

// Whether `val` is in the set (`==`).
func (s *Set[T]) Contains(val T) bool {
	_, exists := s.data[val]
	return exists
}

// Remove `val` if present. Returns whether it was found. Missing: `false`, set unchanged.
func (s *Set[T]) Remove(val T) bool {
	exists := s.Contains(val)
	delete(s.data, val)
	return exists
}

// Drop all elements. Keep backing buckets so later `Add`s can reuse them.
func (s *Set[T]) Clear() {
	clear(s.data)
}

// Copy of the current elements. Order is random. Caller can mutate the result without affecting the set. Empty: non-nil `[]T{}`.
func (s *Set[T]) ToSlice() []T {
	if len(s.data) == 0 {
		return []T{}
	}
	return slices.Collect(maps.Keys(s.data))
}

// Yield elements in random hash order. Use with `for range`. Early `break` stops iteration (`yield` must return `false` and stop).
func (s *Set[T]) All() iter.Seq[T] {
	return maps.Keys(s.data)
}

// Return string representation of set items
func (s *Set[T]) String() string {
	vals := make([]string, 0, len(s.data))
	for v := range s.data {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return "{" + strings.Join(vals, ", ") + "}"
}

// ---- Package level utils ----//
