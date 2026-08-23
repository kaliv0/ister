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
	buf   []T // internal ring buffer
	head  int
	count int // live elements
}

// Empty deque. Zero-capacity backing; first PushFront / PushBack allocates.
func New[T any]() *Deque[T] {
	return &Deque[T]{}
}

// Deque from the given values. First value is the front.
// Empty (no args) is New(). Otherwise one allocation: capacity next power of two
// ≥ max(len(vals), 8), copy into buf[0:n], head = 0. Do not grow in a loop.
func Of[T any](vals ...T) *Deque[T] {
	count := len(vals)
	if count == 0 {
		return New[T]()
	}
	buf := makeBufWithCap[T](count)
	copy(buf, vals)
	return &Deque[T]{buf: buf, head: 0, count: count}
}

// Deque from a slice. Copies s; first element is the front.
// FromSlice(nil) / empty slice is New(). Same capacity rule as Of.
func FromSlice[T any](s []T) *Deque[T] {
	return Of(s...)
}

// Number of elements.
func (d *Deque[T]) Len() int {
	return d.count
}

// Whether the deque has no elements.
func (d *Deque[T]) IsEmpty() bool {
	return d.Len() == 0
}

// Add val at the front. A nil receiver panics — use New or Of.
func (d *Deque[T]) PushFront(val T) {
	count := len(d.buf)
	if count == 0 || count == d.Len() {
		d.grow()
	}

	d.head = d.wrapIdx(-1) // shift head left
	d.buf[d.head] = val
	d.count++
}

// Add val at the back. A nil receiver panics — use New or Of.
func (d *Deque[T]) PushBack(val T) {
	count := len(d.buf)
	if count == 0 || count == d.Len() {
		d.grow()
	}

	idx := d.wrapIdx(d.count) // find "tail"
	d.buf[idx] = val
	d.count++
}

func (d *Deque[T]) grow() {
	newBuf := makeBufWithCap[T](d.count + 1)
	for i := 0; i < d.count; i++ {
		newBuf[i] = d.buf[d.wrapIdx(i)]
	}
	/*
		// faster for large deques when benchmarks ask for it
		d.copyLinear(newBuf)
	*/
	d.buf = newBuf
	d.head = 0
}

// Remove and return the front element. Returns (zero, false) if empty.
func (d *Deque[T]) PopFront() (T, bool) {
	return d.handleHead(true)
}

// Return the front element without removing it. Returns (zero, false) if empty.
func (d *Deque[T]) PeekFront() (T, bool) {
	return d.handleHead(false)
}

func (d *Deque[T]) handleHead(remove bool) (T, bool) {
	var zero T
	if d.IsEmpty() {
		return zero, false
	}

	val := d.buf[d.head]
	if remove {
		d.buf[d.head] = zero
		d.count--
		if d.IsEmpty() {
			d.head = 0
		} else {
			d.head = d.wrapIdx(1) // move head from idx '0' to '1' (wrapped around)
		}
	}
	return val, true
}

// Remove and return the back element. Returns (zero, false) if empty.
func (d *Deque[T]) PopBack() (T, bool) {
	return d.handleTail(true)
}

// Return the back element without removing it. Returns (zero, false) if empty.
func (d *Deque[T]) PeekBack() (T, bool) {
	return d.handleTail(false)
}

func (d *Deque[T]) handleTail(remove bool) (T, bool) {
	var zero T
	if d.IsEmpty() {
		return zero, false
	}

	idx := d.wrapIdx(d.count - 1) // find "tail"
	val := d.buf[idx]
	if remove {
		d.buf[idx] = zero
		d.count--
		if d.IsEmpty() {
			d.head = 0
		}
	}
	return val, true
}

// Drop all elements. Keep backing capacity so later pushes can reuse it.
func (d *Deque[T]) Clear() {
	clear(d.buf)
	d.head = 0
	d.count = 0
}

// Copy of the current elements, front to back. Empty: non-nil []T{}.
func (d *Deque[T]) ToSlice() []T {
	if d.IsEmpty() {
		return []T{}
	}
	return slices.Collect(d.All())
	/*
		// avoid iter overhead on hot paths
		out := make([]T, d.count)
		d.copyLinear(out)
		return out
	*/
}

// Yield elements from front to back. Use with for range. Early break stops iteration.
func (d *Deque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < d.count; i++ {
			idx := d.wrapIdx(i)
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

// Buf index for offset i from head. Mask wrap; capacity must be a power of two.
func (d *Deque[T]) wrapIdx(i int) int {
	idx := (d.head + i) & (len(d.buf) - 1) // (d.head + i) % len(d.buf)
	return idx
}

// capacity for at least count elements; power of two, min 8.
func makeBufWithCap[T any](count int) []T {
	capacity := 1 << bits.Len(uint((max(count, defaultCapacity))-1))
	buf := make([]T, capacity)
	return buf
}

/*
// Copy live elements front-to-back into dst. dst must have length ≥ count.
// One or two copy calls; no wrapIdx loop.
func (d *Deque[T]) copyLinear(dst []T) {
	n := d.count
	if d.head+n <= len(d.buf) {
		copy(dst, d.buf[d.head:d.head+n])
		return
	}
	n1 := len(d.buf) - d.head
	copy(dst, d.buf[d.head:])
	copy(dst[n1:], d.buf[:n-n1])
}
*/
