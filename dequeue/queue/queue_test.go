package queue

import (
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func TestNew(t *testing.T) {
	eqQueue(t, New[int](), []int{})
}

func TestTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		q := Of("a", "b")
		q.Enqueue("c")
		front, ok := q.Peek()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, front, "a")
		val, ok := q.Dequeue()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, "a")
		eqQueue(t, q, []string{"b", "c"})
	})
	t.Run("struct", func(t *testing.T) {
		q := New[testutil.Pair]()
		q.Enqueue(testutil.Pair{A: 1, B: 2})
		q.Enqueue(testutil.Pair{A: 3, B: 4})
		val, ok := q.Dequeue()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, testutil.Pair{A: 1, B: 2})
		eqQueue(t, q, []testutil.Pair{{A: 3, B: 4}})
	})
}

func TestOf(t *testing.T) {
	eqQueue(t, Of(1, 2, 3), []int{1, 2, 3})
	eqQueue(t, Of[int](), []int{})
}

func TestIsEmpty(t *testing.T) {
	testutil.EqVal(t, New[int]().IsEmpty(), true)
	testutil.EqVal(t, Of(1).IsEmpty(), false)
}

func TestOfCopies(t *testing.T) {
	src := make([]int, 3, 8)
	copy(src, []int{1, 2, 3})
	q := Of(src...)
	q.Enqueue(4)
	testutil.EqVal(t, src[:cap(src)][3], 0)
	eqQueue(t, q, []int{1, 2, 3, 4})
}

func TestOfOrder(t *testing.T) {
	q := Of(1, 2, 3)
	front, ok := q.Peek()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, front, 1)
}

func TestFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		q := FromSlice(src)
		src[0] = 99
		eqQueue(t, q, []int{1, 2, 3})
	})
	t.Run("first is front", func(t *testing.T) {
		q := FromSlice([]int{1, 2, 3})
		front, ok := q.Peek()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, front, 1)
	})
	t.Run("nil", func(t *testing.T) {
		q := FromSlice([]int(nil))
		testutil.EqVal(t, q.Len(), 0)
	})
}

func TestZeroValue(t *testing.T) {
	q := &Queue[int]{}
	testutil.EqVal(t, q.Len(), 0)
	q.Enqueue(1)
	eqQueue(t, q, []int{1})
}

func TestLen(t *testing.T) {
	testutil.EqVal(t, New[int]().Len(), 0)
	testutil.EqVal(t, Of(1, 2, 3).Len(), 3)
}

func TestEnqueue(t *testing.T) {
	q := New[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	eqQueue(t, q, []int{1, 2, 3})
}

func TestDequeue(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
		want    []int
	}{
		{"three", []int{1, 2, 3}, 1, true, []int{2, 3}},
		{"one", []int{1}, 1, true, []int{}},
		{"empty", nil, 0, false, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := Of(tc.start...)
			val, ok := q.Dequeue()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			eqQueue(t, q, tc.want)
		})
	}
}

func TestDequeueOrder(t *testing.T) {
	q := Of(1, 2, 3)
	var got []int
	for {
		v, ok := q.Dequeue()
		if !ok {
			break
		}
		got = append(got, v)
	}
	testutil.Eq(t, got, []int{1, 2, 3})
}

func TestAfterDequeue(t *testing.T) {
	q := Of(1, 2, 3)
	q.Dequeue()
	testutil.EqVal(t, q.Len(), 2)
	eqQueue(t, q, []int{2, 3})
	testutil.EqVal(t, q.String(), "[2, 3]")
	var vals []int
	for v := range q.All() {
		vals = append(vals, v)
	}
	testutil.Eq(t, vals, []int{2, 3})
}

func TestDequeueThenEnqueue(t *testing.T) {
	q := Of(1, 2, 3)
	q.Dequeue()
	q.Dequeue()
	q.Enqueue(4)
	eqQueue(t, q, []int{3, 4})
}

func TestDequeueKeepsCapacity(t *testing.T) {
	q := Of(1, 2, 3)
	capBefore := cap(q.items)
	for !q.IsEmpty() {
		q.Dequeue()
	}
	if cap(q.items) != capBefore {
		t.Fatalf("full drain should keep capacity: got %d, want %d", cap(q.items), capBefore)
	}
	q.Enqueue(9)
	eqQueue(t, q, []int{9})
}

func TestDequeueZeros(t *testing.T) {
	a, b := 1, 2
	q := Of(&a, &b)
	q.Dequeue()
	testutil.EqVal(t, q.items[0], (*int)(nil))
}

// compact only when head != 0 and len == cap.
func TestCompact(t *testing.T) {
	t.Run("compacts when full", func(t *testing.T) {
		q := Of(1, 2, 3)
		capBefore := cap(q.items)
		q.Dequeue() // head > 0, len == cap
		q.Enqueue(4)
		testutil.EqVal(t, q.head, 0)
		eqQueue(t, q, []int{2, 3, 4})
		if cap(q.items) != capBefore {
			t.Fatalf("compact should keep capacity: got %d, want %d", cap(q.items), capBefore)
		}
	})
	t.Run("skips when spare cap", func(t *testing.T) {
		q := &Queue[int]{items: make([]int, 3, 8)}
		copy(q.items, []int{1, 2, 3})
		q.Dequeue()
		q.Dequeue() // head == 2, len < cap
		q.Enqueue(4)
		testutil.EqVal(t, q.head, 2)
		eqQueue(t, q, []int{3, 4})
	})
	t.Run("zeros abandoned", func(t *testing.T) {
		a, b, c := 1, 2, 3
		q := Of(&a, &b, &c)
		q.Dequeue() // [nil, &b, &c], head == 1
		q.compact()
		testutil.EqVal(t, q.head, 0)
		eqQueue(t, q, []*int{&b, &c})
		testutil.EqVal(t, q.items[:cap(q.items)][2], (*int)(nil))
	})
}

/*
// Strategy: compact when head != 0 and head >= Len().
func TestCompactA(t *testing.T) {
	q := &Queue[int]{items: make([]int, 4, 10)}
	copy(q.items, []int{1, 2, 3, 4})
	q.Dequeue()
	q.Dequeue() // head == 2, Len == 2, spare cap (D would skip)
	q.Enqueue(5)
	testutil.EqVal(t, q.head, 0)
	eqQueue(t, q, []int{3, 4, 5})
}
*/

func TestPeek(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
	}{
		{"three", []int{1, 2, 3}, 1, true},
		{"one", []int{1}, 1, true},
		{"empty", nil, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := Of(tc.start...)
			val, ok := q.Peek()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			testutil.EqVal(t, q.Len(), len(tc.start))
		})
	}
}

func TestClear(t *testing.T) {
	q := Of(1, 2, 3, 4, 5)
	capBefore := cap(q.items)
	q.Clear()
	testutil.EqVal(t, q.Len(), 0)
	if cap(q.items) != capBefore {
		t.Fatalf("Clear should keep capacity: got %d, want %d", cap(q.items), capBefore)
	}
	q.Enqueue(9)
	eqQueue(t, q, []int{9})
}

func TestClearResetsHead(t *testing.T) {
	q := Of(1, 2, 3)
	q.Dequeue()
	q.Dequeue()
	q.Clear()
	q.Enqueue(9)
	front, ok := q.Peek()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, front, 9)
}

func TestAll(t *testing.T) {
	t.Run("front to back", func(t *testing.T) {
		var vals []int
		for v := range Of(1, 2, 3).All() {
			vals = append(vals, v)
		}
		testutil.Eq(t, vals, []int{1, 2, 3})
	})
	t.Run("early break", func(t *testing.T) {
		var vals []int
		for v := range Of(1, 2, 3).All() {
			vals = append(vals, v)
			break
		}
		testutil.Eq(t, vals, []int{1})
	})
	t.Run("empty", func(t *testing.T) {
		for range New[int]().All() {
			t.Fatal("empty queue should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	q := Of(1, 2, 3)
	a := q.ToSlice()
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	testutil.Eq(t, q.ToSlice(), []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := Of(1, 2)
	cleared.Clear()
	drained := Of(1, 2)
	drained.Dequeue()
	drained.Dequeue()
	cases := []struct {
		name string
		got  []int
	}{
		{"New", New[int]().ToSlice()},
		{"Of", Of[int]().ToSlice()},
		{"FromSlice nil", FromSlice([]int(nil)).ToSlice()},
		{"after Clear", cleared.ToSlice()},
		{"after drain", drained.ToSlice()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.NonNilEmpty(t, tc.got)
		})
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ints", Of(1, 2, 3).String(), "[1, 2, 3]"},
		{"empty", New[int]().String(), "[]"},
		{"stringer", Of(testutil.Ident(1), testutil.Ident(2)).String(), "[id=1, id=2]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.got, tc.want)
		})
	}
}

func TestNilReceiverPanics(t *testing.T) {
	var q *Queue[int]
	cases := []struct {
		name string
		fn   func()
	}{
		{"Enqueue", func() { q.Enqueue(1) }},
		{"Len", func() { q.Len() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func eqQueue[T comparable](t *testing.T, q *Queue[T], want []T) {
	t.Helper()
	testutil.Eq(t, q.ToSlice(), want)
}
