package list

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"
)

type List[T comparable] struct {
	items []T
}

// Empty list. Zero-capacity backing slice; first Add allocates.
func New[T comparable]() *List[T] {
	return &List[T]{}
}

// List from the given values, in order. Of[T]() with no args is empty.
func Of[T comparable](vals ...T) *List[T] {
	if len(vals) == 0 {
		return New[T]()
	}
	return &List[T]{items: slices.Clone(vals)}
}

// List from a slice. Copy `s` so later mutations of the original slice do not affect the list. `FromSlice(nil)` is empty.
func FromSlice[T comparable](s []T) *List[T] {
	return Of(s...)
}

// Number of elements.
func (l *List[T]) Len() int {
	return len(l.items)
}

// Add one or more elements. A nil receiver panics — use `New` or `Of`.
func (l *List[T]) Add(vals ...T) {
	l.items = append(l.items, vals...)
}

// Insert at index `i`, shifting later elements right. `i == Len()` is append.
func (l *List[T]) Insert(i int, val ...T) bool {
	if i < 0 || i > l.Len() {
		return false
	}
	l.items = slices.Insert(l.items, i, val...)
	return true
}

// Element at index i.
func (l *List[T]) Get(i int) (T, bool) {
	if l.inRange(i) {
		var zero T
		return zero, false
	}
	return l.items[i], true
}

// Replace the element at `i`.
func (l *List[T]) Set(i int, val T) bool {
	if l.inRange(i) {
		return false
	}
	l.items[i] = val
	return true
}

// Remove the element at `i`, shift later elements left, return the removed value.
func (l *List[T]) RemoveAt(i int) (T, bool) {
	if l.inRange(i) {
		var zero T
		return zero, false
	}
	val := l.items[i]
	l.items = slices.Delete(l.items, i, i+1)
	return val, true
}

func (l *List[T]) inRange(i int) bool {
	return i < 0 || i >= l.Len()
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

// Remove the first element equal to `val`. Returns whether it was found. Missing: `false`, list unchanged. A nil receiver panics — use `New` or `Of`.
func (l *List[T]) Remove(val T) bool {
	i := slices.Index(l.items, val)
	if i == -1 {
		return false
	}

	l.items = slices.Delete(l.items, i, i+1)
	return true
}

// Index of the first element equal to `val`, or `-1`. A nil receiver panics — use `New` or `Of`.
func (l *List[T]) Index(val T) int {
	return slices.Index(l.items, val)
}

// Whether `val` is in the list (`==`). A nil receiver panics — use `New` or `Of`.
func (l *List[T]) Contains(val T) bool {
	return slices.Contains(l.items, val)
}

// Copy of the current elements. Caller can mutate the result without affecting the list.
func (l *List[T]) ToSlice() []T {
	if l.Len() == 0 {
		return []T{}
	}
	return slices.Clone(l.items)
}

// Yield elements from front to back. Use with `for range`. Early `break` stops iteration.
func (l *List[T]) All() iter.Seq[T] {
	return slices.Values(l.items)
}

// String representation of list items
func (l *List[T]) String() string {
	vals := make([]string, 0, l.Len())
	for _, v := range l.items {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

// Sort all elements. Stays package-level because T must be cmp.Ordered, stricter than List's comparable. A nil `l` panics — use `New` or `Of`.
func Sort[T cmp.Ordered](l *List[T]) {
	slices.Sort(l.items)
}
