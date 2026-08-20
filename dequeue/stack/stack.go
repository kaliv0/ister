package stack

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

type Stack[T any] struct {
	items []T
}

// Empty stack. Zero-capacity backing slice; first Push allocates.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Stack from the given values. Last value is the top.
func Of[T any](vals ...T) *Stack[T] {
	if len(vals) == 0 {
		return New[T]()
	}
	return &Stack[T]{items: slices.Clone(vals)}
}

// Stack from a slice. Copies s; last element is the top. FromSlice(nil) is empty.
func FromSlice[T any](s []T) *Stack[T] {
	return Of(s...)
}

// Number of elements.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Whether the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Add val to the top. A nil receiver panics — use `New` or `Of`.
func (s *Stack[T]) Push(val T) {
	s.items = append(s.items, val)
}

// Remove and return the top element. Returns (zero, false) if empty.
func (s *Stack[T]) Pop() (T, bool) {
	return s.handleLast(true)
}

// Return the top element without removing it. Returns (zero, false) if empty.
func (s *Stack[T]) Peek() (T, bool) {
	return s.handleLast(false)
}

func (s *Stack[T]) handleLast(remove bool) (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}

	val := s.items[s.Len()-1]
	if remove {
		s.items = slices.Delete(s.items, s.Len()-1, s.Len())
	}
	return val, true
}

// Drop all elements. Keep backing capacity so later Pushes can reuse it.
func (s *Stack[T]) Clear() {
	clear(s.items)
	s.items = s.items[:0]
}

// Copy of the current elements, bottom to top. Empty: non-nil []T{}.
func (s *Stack[T]) ToSlice() []T {
	if s.IsEmpty() {
		return []T{}
	}
	return slices.Clone(s.items)
}

// Yield elements from top to bottom (LIFO order). Use with for range. Early break stops iteration.
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range slices.Backward(s.items) {
			if !yield(v) {
				return
			}
		}
	}
}

// String representation, top to bottom.
func (s *Stack[T]) String() string {
	vals := make([]string, 0, s.Len())
	for v := range s.All() {
		vals = append(vals, fmt.Sprintf("%v", v))
	}

	// TODO: add type (e.g. "stack[1, 2, 3]") for all String funcs -> extract common util
	return "[" + strings.Join(vals, ", ") + "]"
}
