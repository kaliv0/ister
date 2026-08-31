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
func New[T comparable]() *Set[T] {
	return &Set[T]{}
}

// Set from the given values. Duplicates are dropped. Of[T]() with no args is empty.
func Of[T comparable](vals ...T) *Set[T] {
	s := New[T]()
	s.Add(vals...)
	return s
}

// Set from a slice. Copy `s` so later mutations of the original slice do not affect the set.
// Duplicates are dropped. `FromSlice(nil)` is empty.
func FromSlice[T comparable](s []T) *Set[T] {
	return Of(s...)
}

// Number of elements.
func (s *Set[T]) Len() int {
	return len(s.data)
}

// Add one or more elements. A nil receiver panics — use `New` or `Of`.
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
	if s.Len() == 0 {
		return []T{}
	}
	return slices.Collect(maps.Keys(s.data))
}

// Yield elements in random hash order. Use with `for range`. Early `break` stops iteration (`yield` must return `false` and stop).
func (s *Set[T]) All() iter.Seq[T] {
	return maps.Keys(s.data)
}

// String representation of set items
func (s *Set[T]) String() string {
	vals := make([]string, 0, s.Len())
	for v := range s.data {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return "{" + strings.Join(vals, ", ") + "}"
}

//---- Algebra ----//

// Elements in `s` or `other` (or both). A nil `other` is treated as empty.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	other = coalesce(other)
	result := New[T]()
	if n := s.Len() + other.Len(); n > 0 {
		result.data = make(map[T]struct{}, n)
		maps.Copy(result.data, s.data)
		maps.Copy(result.data, other.data)
	}
	return result
}

// Elements in both `s` and `other`. A nil `other` is treated as empty.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	other = coalesce(other)
	result := New[T]()
	if n := s.Len() + other.Len(); n == 0 {
		return result
	}

	small, big := compare(s, other)
	result.data = make(map[T]struct{}, small.Len())

	for v := range small.All() {
		if big.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Elements in `s` but not in `other`. A nil `other` is treated as empty.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	other = coalesce(other)
	result := New[T]()
	if n := s.Len() + other.Len(); n > 0 {
		result.data = make(map[T]struct{}, s.Len())

		for v := range s.All() {
			if !other.Contains(v) {
				result.Add(v)
			}
		}
	}
	return result
}

// Elements in `s` or `other` but not both. A nil `other` is treated as empty.
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	other = coalesce(other)
	result := New[T]()
	if n := s.Len() + other.Len(); n > 0 {
		result.data = make(map[T]struct{}, n)

		for v := range s.All() {
			if !other.Contains(v) {
				result.Add(v)
			}
		}

		for v := range other.All() {
			if !s.Contains(v) {
				result.Add(v)
			}
		}
	}
	return result
}

// Whether `s` and `other` contain the same elements. Order does not matter. A nil `other` is treated as empty.
func (s *Set[T]) Equal(other *Set[T]) bool {
	other = coalesce(other)
	return maps.Equal(s.data, other.data)
}

// Whether every element of `s` is in `other`. Empty set is a subset of every set. A nil `other` is treated as empty.
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
	return isSubset(s, other)
}

// Whether every element of `other` is in `s`. Convenience inverse of IsSubsetOf. A nil `other` is treated as empty.
func (s *Set[T]) IsSupersetOf(other *Set[T]) bool {
	return isSubset(other, s)
}

func isSubset[T comparable](a, b *Set[T]) bool {
	a, b = coalesce(a), coalesce(b)
	if a.Len() == 0 {
		return true
	}
	if b.Len() < a.Len() {
		return false
	}

	for v := range a.All() {
		if !b.Contains(v) {
			return false
		}
	}
	return true
}

// Whether `s` and `other` share no elements. A nil `other` is treated as empty.
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	other = coalesce(other)
	if s.Len()+other.Len() == 0 {
		return true
	}

	small, big := compare(s, other)
	for v := range small.All() {
		if big.Contains(v) {
			return false
		}
	}
	return true
}

func compare[T comparable](a, b *Set[T]) (small, big *Set[T]) {
	if a.Len() < b.Len() {
		return a, b
	}
	return b, a
}

func coalesce[T comparable](s *Set[T]) *Set[T] {
	if s != nil {
		return s
	}
	return New[T]()
}
