package queue

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

type Queue[T any] struct {
	items []T
	head  int
}

// Empty queue. Zero-capacity backing slice; first Enqueue allocates.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Queue from the given values. First value is the front.
func Of[T any](vals ...T) *Queue[T] {
	if len(vals) == 0 {
		return New[T]()
	}
	return &Queue[T]{items: slices.Clone(vals)}
}

// Queue from a slice. Copies s; first element is the front. FromSlice(nil) is empty.
func FromSlice[T any](s []T) *Queue[T] {
	return Of(s...)
}

// Number of elements. Live range is items[head:].
func (q *Queue[T]) Len() int {
	return len(q.items) - q.head
}

// Whether the queue has no elements.
func (q *Queue[T]) IsEmpty() bool {
	return q.Len() == 0
}

// Add val at the back. A nil receiver panics — use New or Of.
func (q *Queue[T]) Enqueue(val T) {
	/*
		// waste > live
		if q.head != 0 && q.head >= q.Len() {
			q.compact()
		}
	*/

	// cap is full - avoid realloc
	if q.head != 0 && len(q.items) == cap(q.items) {
		q.compact()
	}
	q.items = append(q.items, val)
}

func (q *Queue[T]) compact() {
	n := copy(q.items, q.view())
	clear(q.items[n:])
	q.reset(n)
}

// Remove and return the front element. Returns (zero, false) if empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	return q.handleHead(true)
}

// Return the front element without removing it. Returns (zero, false) if empty.
func (q *Queue[T]) Peek() (T, bool) {
	return q.handleHead(false)
}

func (q *Queue[T]) handleHead(remove bool) (T, bool) {
	var zero T
	if q.IsEmpty() {
		return zero, false
	}

	val := q.items[q.head]
	if remove {
		q.items[q.head] = zero
		q.head++
		if q.head == len(q.items) {
			q.reset(0)
		}
	}
	return val, true
}

// Drop all elements. Keep backing capacity so later Enqueues can reuse it.
func (q *Queue[T]) Clear() {
	clear(q.items)
	q.reset(0)
}

func (q *Queue[T]) reset(n int) {
	q.items = q.items[:n]
	q.head = 0
}

// Copy of the current elements, front to back. Empty: non-nil []T{}.
func (q *Queue[T]) ToSlice() []T {
	if q.IsEmpty() {
		return []T{}
	}
	return slices.Clone(q.view())
}

// Yield elements from front to back. Use with for range. Early break stops iteration.
func (q *Queue[T]) All() iter.Seq[T] {
	return slices.Values(q.view())
}

// String representation, front to back.
func (q *Queue[T]) String() string {
	vals := make([]string, 0, q.Len())
	for _, v := range q.view() {
		vals = append(vals, fmt.Sprintf("%v", v))
	}

	return "[" + strings.Join(vals, ", ") + "]"
}

func (q *Queue[T]) view() []T {
	return q.items[q.head:]
}
