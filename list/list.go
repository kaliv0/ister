package list

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"
)

type List[T any] struct {
	items []T
}

// Empty list. Zero-capacity backing slice; first Add allocates.
func New[T any]() *List[T] {
	return &List[T]{items: []T{}}
}

// List from the given values, in order. Of[T]() with no args is empty.
func Of[T any](vals ...T) *List[T] {
	return &List[T]{items: slices.Clone(vals)}
}

// List from a slice. Copy `s` so later mutations of the original slice do not affect the list. `FromSlice(nil)` is empty.
func FromSlice[T any](s []T) *List[T] {
	return &List[T]{items: slices.Clone(s)}
}

// Number of elements.
func (l *List[T]) Len() int {
	return len(l.items)
}

func (l *List[T]) Add(vals ...T) {
	l.items = append(l.items, vals...)
}

// Insert at index `i`, shifting later elements right. `i == Len()` is append. Panics if `i < 0` or `i > Len()`.
func (l *List[T]) Insert(i int, val ...T) {
	l.items = slices.Insert(l.items, i, val...)
}

// Element at index i. Panics if out of range.
func (l *List[T]) Get(i int) T {
	return l.items[i]
}

// Replace the element at `i`. Panics if out of range.
func (l *List[T]) Set(i int, val T) {
	l.items[i] = val
}

// Remove the element at `i`, shift later elements left, return the removed value. Panics if out of range.
func (l *List[T]) RemoveAt(i int) T {
	val := l.items[i]
	l.items = slices.Delete(l.items, i, i+1)
	return val
}

// Drop all elements. Keep backing capacity so later `Add`s can reuse it.
func (l *List[T]) Clear() {
	clear(l.items)
	l.items = l.items[:0]
}

// Reverse elements.
func (l *List[T]) Reverse() {
	slices.Reverse(l.items)
}

// Copy of the current elements. Caller can mutate the result without affecting the list.
func (l *List[T]) ToSlice() []T {
	if len(l.items) == 0 {
		return []T{}
	}
	return slices.Clone(l.items)
}

// Yield elements from front to back. Use with `for range`. Early `break` stops iteration (`yield` must return `false` and stop).
func (l *List[T]) All() iter.Seq[T] {
	return slices.Values(l.items)
}

// Return string representation of list items
func (l *List[T]) String() string {
	vals := make([]string, 0, len(l.items))
	for _, v := range l.items {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

// ---- Package level utils ----//
// Remove the first element equal to `val`. Returns whether it was found. Missing: `false`, list unchanged.
func Remove[T comparable](l *List[T], val T) bool {
	i := slices.Index(l.items, val)
	if i == -1 {
		return false
	}

	l.items = slices.Delete(l.items, i, i+1)
	return true
}

// Index of the first element equal to `val`, or `-1`.
func Index[T comparable](l *List[T], val T) int {
	return slices.Index(l.items, val)
}

// Whether `val` is in the list (`==`).
func Contains[T comparable](l *List[T], val T) bool {
	return slices.Contains(l.items, val)
}

// Sort all elements.
func Sort[T cmp.Ordered](l *List[T]) {
	slices.Sort(l.items)
}
