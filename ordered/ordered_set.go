package ordered

import (
	"fmt"
	"iter"
	"strings"

	"github.com/kaliv0/ister/internal/rbtree"
)

// OrderedSet is a sorted set of unique elements.
// Use New or Of; the zero value has a nil tree and panics on use.
type OrderedSet[T comparable] struct {
	tree *rbtree.Tree[T, struct{}]
}

// Empty set. Wraps an empty rbtree (nil root until first Add).
// A nil less panics when ordering is needed, not in New.
// Use New or Of; the zero value has a nil tree and panics on use.
func New[T comparable](less func(a, b T) bool) *OrderedSet[T] {
	return &OrderedSet[T]{tree: rbtree.New[T, struct{}](less)}
}

// Set from the given values, sorted unique. Duplicates dropped by less.
// Empty (no args) is New(less). Nil less: same deferred panic as PriorityQueue.
func Of[T comparable](less func(a, b T) bool, vals ...T) *OrderedSet[T] {
	s := New(less)
	s.Add(vals...)
	return s
}

// Set from a slice. Copies values into the set. FromSlice(nil) is empty.
func FromSlice[T comparable](less func(a, b T) bool, s []T) *OrderedSet[T] {
	return Of(less, s...)
}

// Number of elements.
func (s *OrderedSet[T]) Len() int {
	return s.tree.Len()
}

// Add one or more elements. Duplicates (by less) are ignored.
// A nil receiver panics — use New or Of.
func (s *OrderedSet[T]) Add(vals ...T) {
	for _, v := range vals {
		s.tree.Put(v, struct{}{})
	}
}

// Whether val is in the set (by less equality).
func (s *OrderedSet[T]) Contains(val T) bool {
	return s.tree.Contains(val)
}

// Remove val if present. Returns whether it was found. Missing: false, set unchanged.
func (s *OrderedSet[T]) Remove(val T) bool {
	_, exists := s.tree.Delete(val)
	return exists
}

// Drop all elements. Keeps less so later Adds can reuse the set.
func (s *OrderedSet[T]) Clear() {
	s.tree.Clear()
}

// Copy of elements in ascending order. Empty: non-nil []T{}.
func (s *OrderedSet[T]) ToSlice() []T {
	if s.Len() == 0 {
		return []T{}
	}
	result := make([]T, 0, s.Len())
	for k := range s.tree.Ascend() {
		result = append(result, k)
	}
	return result
}

// Yield elements in ascending order. Early break stops iteration.
// Do not mutate the set during iteration.
func (s *OrderedSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.tree.Ascend() {
			if !yield(v) {
				return
			}
		}
	}
}

// String representation in ascending order.
func (s *OrderedSet[T]) String() string {
	vals := make([]string, 0, s.Len())
	for v := range s.All() {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return "{" + strings.Join(vals, ", ") + "}"
}

//---- Navigable ----//

// Smallest element. Empty: (zero, false).
func (s *OrderedSet[T]) First() (T, bool) {
	panic("not implemented")
}

// Largest element. Empty: (zero, false).
func (s *OrderedSet[T]) Last() (T, bool) {
	panic("not implemented")
}

// Greatest element ≤ val. None: (zero, false).
func (s *OrderedSet[T]) Floor(val T) (T, bool) {
	panic("not implemented")
}

// Least element ≥ val. None: (zero, false).
func (s *OrderedSet[T]) Ceiling(val T) (T, bool) {
	panic("not implemented")
}

// Greatest element < val. None: (zero, false).
func (s *OrderedSet[T]) Lower(val T) (T, bool) {
	panic("not implemented")
}

// Least element > val. None: (zero, false).
func (s *OrderedSet[T]) Higher(val T) (T, bool) {
	panic("not implemented")
}

// Ascending elements in [lo, hi] (inclusive). Empty range yields nothing.
func (s *OrderedSet[T]) Range(lo, hi T) iter.Seq[T] {
	panic("not implemented")
}

//---- Algebra ----//

// Elements in s or other (or both). Nil other is empty. Result uses the receiver's less.
func (s *OrderedSet[T]) Union(other *OrderedSet[T]) *OrderedSet[T] {
	panic("not implemented")
}

// Elements in both s and other. Nil other is empty. Result uses the receiver's less.
func (s *OrderedSet[T]) Intersect(other *OrderedSet[T]) *OrderedSet[T] {
	panic("not implemented")
}

// Elements in s but not in other. Nil other is empty. Result uses the receiver's less.
func (s *OrderedSet[T]) Difference(other *OrderedSet[T]) *OrderedSet[T] {
	panic("not implemented")
}

// Elements in s or other but not both. Nil other is empty. Result uses the receiver's less.
func (s *OrderedSet[T]) SymmetricDifference(other *OrderedSet[T]) *OrderedSet[T] {
	panic("not implemented")
}

// Whether s and other contain the same elements (by less / membership). Order does not matter for equality of sets. Nil other is empty.
func (s *OrderedSet[T]) Equal(other *OrderedSet[T]) bool {
	panic("not implemented")
}

// Whether every element of s is in other. Empty is a subset of every set. Nil other is empty.
func (s *OrderedSet[T]) IsSubsetOf(other *OrderedSet[T]) bool {
	panic("not implemented")
}

// Whether every element of other is in s. Nil other is empty.
func (s *OrderedSet[T]) IsSupersetOf(other *OrderedSet[T]) bool {
	panic("not implemented")
}

// Whether s and other share no elements. Nil other is empty.
func (s *OrderedSet[T]) IsDisjoint(other *OrderedSet[T]) bool {
	panic("not implemented")
}
