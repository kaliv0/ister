package heap

import (
	stdheap "container/heap"
	"fmt"
	"iter"
	"slices"
	"strings"
)

type PriorityQueue[T comparable] struct {
	h pqHeap[T]
}

type pqHeap[T comparable] struct {
	items []T
	less  func(a, b T) bool
}

// Empty queue. Zero-capacity backing slice; first Push allocates.
// A nil less panics when ordering is needed (Push, heapify), not in New.
func New[T comparable](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		h: pqHeap[T]{less: less},
	}
}

// Queue from the given values, heapified. Empty (no args) is New(less).
// A nil less panics when heapify runs (non-empty vals); empty Of does not.
func Of[T comparable](less func(a, b T) bool, vals ...T) *PriorityQueue[T] {
	if len(vals) == 0 {
		return New[T](less)
	}
	h := pqHeap[T]{items: slices.Clone(vals), less: less}
	stdheap.Init(&h)
	pq := &PriorityQueue[T]{h}
	return pq
}

// Queue from a slice. Copies s, then heapifies. FromSlice(nil) is empty.
func FromSlice[T comparable](less func(a, b T) bool, s []T) *PriorityQueue[T] {
	return Of(less, s...)
}

// Number of elements.
func (pq *PriorityQueue[T]) Len() int {
	return pq.h.Len()
}

// Whether the queue has no elements.
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Len() == 0
}

// Insert val. O(log n). A nil receiver panics — use New or Of.
func (pq *PriorityQueue[T]) Push(val T) {
	stdheap.Push(&pq.h, val)
}

// Remove and return the highest-priority element. Returns (zero, false) if empty.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	return pq.handleHead(true)
}

// Return the highest-priority element without removing it. Returns (zero, false) if empty.
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	return pq.handleHead(false)
}

func (pq *PriorityQueue[T]) handleHead(remove bool) (T, bool) {
	if pq.IsEmpty() {
		var zero T
		return zero, false
	}

	var val T
	if remove {
		val = stdheap.Pop(&pq.h).(T)
	} else {
		val = pq.h.items[0]
	}
	return val, true
}

// Remove the first element equal to val (==).
// O(n) linear scan to find the element, then O(log n) stdheap fix via stdheap.Remove.
// Returns whether val was found.
func (pq *PriorityQueue[T]) Remove(val T) bool {
	i := slices.Index(pq.h.items, val)
	if i == -1 {
		return false
	}
	stdheap.Remove(&pq.h, i)
	return true
}

// Drop all elements. Keep backing capacity and less so later Pushes can reuse the queue.
func (pq *PriorityQueue[T]) Clear() {
	clear(pq.h.items)
	pq.h.items = pq.h.items[:0]
}

// Copy of the current elements in heap layout order (not poll order, not sorted).
// Empty: non-nil []T{}.
func (pq *PriorityQueue[T]) ToSlice() []T {
	if pq.IsEmpty() {
		return []T{}
	}
	return slices.Clone(pq.h.items)
}

// Yield elements in heap layout order. Use with for range. Early break stops iteration.
// Not poll order. Do not mutate the queue during iteration.
func (pq *PriorityQueue[T]) All() iter.Seq[T] {
	return slices.Values(pq.h.items)
}

// String representation in heap layout order.
func (pq *PriorityQueue[T]) String() string {
	vals := make([]string, 0, pq.Len())
	for _, v := range pq.h.items {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

// ----- heap.Interface ----- //
func (h *pqHeap[T]) Len() int {
	return len(h.items)
}

func (h *pqHeap[T]) Less(i, j int) bool {
	return h.less(h.items[i], h.items[j])
}

func (h *pqHeap[T]) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *pqHeap[T]) Push(x any) {
	h.items = append(h.items, x.(T))
}

func (h *pqHeap[T]) Pop() any {
	n := h.Len() - 1
	val := h.items[n]
	// drop refs for GC
	var zero T
	h.items[n] = zero
	h.items = h.items[:n]

	return val
}
