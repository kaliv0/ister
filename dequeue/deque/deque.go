package deque

import (
	"fmt"
	"iter"
	"math/bits"
	"slices"
	"strings"
)

const defaultCapacity = 8

type Deque[T any] struct {
	buf  []T
	head int
	n    int // live count
}

// Empty deque. Zero-capacity backing; first PushFront / PushBack allocates.
func New[T any]() *Deque[T] {
	return &Deque[T]{}
}

// Deque from the given values. First value is the front.
// Empty (no args) is New(). Otherwise one allocation: capacity next power of two
// ≥ max(len(vals), 8), copy into buf[0:n], head = 0. Do not grow in a loop.
func Of[T any](vals ...T) *Deque[T] {
	if len(vals) == 0 {
		return New[T]()
	}
	count := len(vals)
	capacity := 1 << bits.Len(uint((max(count, defaultCapacity))-1))
	buf := make([]T, capacity)
	copy(buf, vals)
	return &Deque[T]{buf: buf, head: 0, n: count}
}

// Deque from a slice. Copies s; first element is the front.
// FromSlice(nil) / empty slice is New(). Same capacity rule as Of.
func FromSlice[T any](s []T) *Deque[T] {
	return Of(s...)
}

// Number of elements.
func (d *Deque[T]) Len() int {
	return d.n
}

// Whether the deque has no elements.
func (d *Deque[T]) IsEmpty() bool {
	return d.Len() == 0
}

// Drop all elements. Keep backing capacity so later pushes can reuse it.
func (d *Deque[T]) Clear() {
	clear(d.buf)
	d.head = 0
	d.n = 0
}

// Copy of the current elements, front to back. Empty: non-nil []T{}.
func (d *Deque[T]) ToSlice() []T {
	if d.IsEmpty() {
		return []T{}
	}
	return slices.Collect(d.All())
}

// Yield elements from front to back. Use with for range. Early break stops iteration.
func (d *Deque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < d.n; i++ {
			idx := (d.head + i) & (len(d.buf) - 1) // equivalent to (d.head + i) % len(d.buf)
			val := d.buf[idx]
			if !yield(val) {
				return
			}
		}
	}
}

// String representation, front to back.
func (d *Deque[T]) String() string {
	vals := make([]string, 0, d.Len())
	for v := range d.All() {
		vals = append(vals, fmt.Sprintf("%v", v))
	}

	return "[" + strings.Join(vals, ", ") + "]"
}
