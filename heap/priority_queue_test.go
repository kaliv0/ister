package heap

import (
	"slices"
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func lessInt(a, b int) bool { return a < b }

func greaterInt(a, b int) bool { return a > b }

func lessPair(a, b testutil.Pair) bool {
	if a.A != b.A {
		return a.A < b.A
	}
	return a.B < b.B
}

func TestNew(t *testing.T) {
	pq := New(lessInt)
	testutil.EqVal(t, pq.Len(), 0)
	testutil.EqVal(t, pq.IsEmpty(), true)
}

func TestTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		less := func(a, b string) bool { return a < b }
		pq := Of(less, "b", "a")
		pq.Push("c")
		top, ok := pq.Peek()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, top, "a")
		val, ok := pq.Pop()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, "a")
		eqPoll(t, pq, []string{"b", "c"})
	})
	t.Run("struct", func(t *testing.T) {
		pq := New(lessPair)
		pq.Push(testutil.Pair{A: 2, B: 1})
		pq.Push(testutil.Pair{A: 1, B: 9})
		val, ok := pq.Pop()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, testutil.Pair{A: 1, B: 9})
		eqPoll(t, pq, []testutil.Pair{{A: 2, B: 1}})
	})
}

func TestOf(t *testing.T) {
	eqElems(t, Of(lessInt, 3, 1, 2), []int{1, 2, 3})
	eqElems(t, Of(lessInt), []int{})
}

func TestIsEmpty(t *testing.T) {
	testutil.EqVal(t, New(lessInt).IsEmpty(), true)
	testutil.EqVal(t, Of(lessInt, 1).IsEmpty(), false)
}

func TestOfCopies(t *testing.T) {
	src := make([]int, 3, 8)
	copy(src, []int{1, 2, 3})
	pq := Of(lessInt, src...)
	pq.Push(4)
	testutil.EqVal(t, src[:cap(src)][3], 0)
	eqElems(t, pq, []int{1, 2, 3, 4})
}

func TestFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{3, 1, 2}
		pq := FromSlice(lessInt, src)
		src[0] = 99
		eqElems(t, pq, []int{1, 2, 3})
	})
	t.Run("nil", func(t *testing.T) {
		pq := FromSlice(lessInt, []int(nil))
		testutil.EqVal(t, pq.Len(), 0)
	})
	t.Run("empty", func(t *testing.T) {
		pq := FromSlice(lessInt, []int{})
		testutil.EqVal(t, pq.Len(), 0)
		testutil.NonNilEmpty(t, pq.ToSlice())
	})
}

func TestZeroValue(t *testing.T) {
	pq := &PriorityQueue[int]{h: pqHeap[int]{less: lessInt}}
	testutil.EqVal(t, pq.Len(), 0)
	testutil.EqVal(t, pq.IsEmpty(), true)
	pq.Push(1)
	eqElems(t, pq, []int{1})
}

func TestLen(t *testing.T) {
	testutil.EqVal(t, New(lessInt).Len(), 0)
	testutil.EqVal(t, Of(lessInt, 1, 2, 3).Len(), 3)
}

func TestPush(t *testing.T) {
	pq := New(lessInt)
	pq.Push(3)
	pq.Push(1)
	pq.Push(2)
	eqPoll(t, pq, []int{1, 2, 3})
}

func TestPop(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
		want    []int
	}{
		{"three", []int{3, 1, 2}, 1, true, []int{2, 3}},
		{"one", []int{1}, 1, true, []int{}},
		{"empty", nil, 0, false, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pq := Of(lessInt, tc.start...)
			val, ok := pq.Pop()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			eqPoll(t, pq, tc.want)
		})
	}
}

func TestPopOrder(t *testing.T) {
	t.Run("min", func(t *testing.T) {
		eqPoll(t, Of(lessInt, 5, 1, 4, 2, 3), []int{1, 2, 3, 4, 5})
	})
	t.Run("max", func(t *testing.T) {
		eqPoll(t, Of(greaterInt, 5, 1, 4, 2, 3), []int{5, 4, 3, 2, 1})
	})
	t.Run("duplicates", func(t *testing.T) {
		eqPoll(t, Of(lessInt, 2, 1, 2), []int{1, 2, 2})
	})
}

func TestPushAfterPop(t *testing.T) {
	pq := Of(lessInt, 3, 1, 2)
	val, ok := pq.Pop()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, val, 1)
	pq.Push(0)
	eqPoll(t, pq, []int{0, 2, 3})
}

func TestPeek(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
	}{
		{"three", []int{3, 1, 2}, 1, true},
		{"one", []int{1}, 1, true},
		{"empty", nil, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pq := Of(lessInt, tc.start...)
			val, ok := pq.Peek()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			testutil.EqVal(t, pq.Len(), len(tc.start))
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name  string
		start []int
		val   int
		ok    bool
		want  []int
	}{
		{"present", []int{3, 1, 2}, 2, true, []int{1, 3}},
		{"missing", []int{1, 2}, 9, false, []int{1, 2}},
		{"root", []int{3, 1, 2}, 1, true, []int{2, 3}},
		{"only", []int{1}, 1, true, []int{}},
		{"empty", nil, 1, false, []int{}},
		{"duplicate first", []int{2, 1, 2}, 2, true, []int{1, 2}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pq := Of(lessInt, tc.start...)
			testutil.EqVal(t, pq.Remove(tc.val), tc.ok)
			eqPoll(t, pq, tc.want)
		})
	}
}

func TestClear(t *testing.T) {
	pq := Of(lessInt, 1, 2, 3, 4, 5)
	capBefore := cap(pq.h.items)
	pq.Clear()
	testutil.EqVal(t, pq.Len(), 0)
	testutil.EqVal(t, pq.IsEmpty(), true)
	if cap(pq.h.items) != capBefore {
		t.Fatalf("Clear should keep capacity: got %d, want %d", cap(pq.h.items), capBefore)
	}
	// less must still work: second Push compares.
	pq.Push(9)
	pq.Push(3)
	eqElems(t, pq, []int{3, 9})
	got, ok := pq.Peek()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, got, 3)
}

func TestAll(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		pq := Of(lessInt, 3, 1, 2)
		var vals []int
		for v := range pq.All() {
			vals = append(vals, v)
		}
		testutil.Eq(t, vals, pq.ToSlice())
		slices.Sort(vals)
		testutil.Eq(t, vals, []int{1, 2, 3})
	})
	t.Run("early break", func(t *testing.T) {
		n := 0
		for range Of(lessInt, 3, 1, 2).All() {
			n++
			break
		}
		testutil.EqVal(t, n, 1)
	})
	t.Run("empty", func(t *testing.T) {
		for range New(lessInt).All() {
			t.Fatal("empty queue should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	pq := Of(lessInt, 3, 1, 2)
	a := pq.ToSlice()
	slices.Sort(a)
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	got := pq.ToSlice()
	slices.Sort(got)
	testutil.Eq(t, got, []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := Of(lessInt, 1, 2)
	cleared.Clear()
	drained := Of(lessInt, 1, 2)
	drained.Pop()
	drained.Pop()
	cases := []struct {
		name string
		got  []int
	}{
		{"New", New(lessInt).ToSlice()},
		{"Of", Of(lessInt).ToSlice()},
		{"FromSlice nil", FromSlice(lessInt, []int(nil)).ToSlice()},
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
		{"empty", New(lessInt).String(), "[]"},
		{"single", Of(lessInt, 1).String(), "[1]"},
		{"stringer", Of(func(a, b testutil.Ident) bool { return a < b }, testutil.Ident(1)).String(), "[id=1]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.got, tc.want)
		})
	}
}

// Nil less is deferred: empty New/Of succeed; panic on first compare (needs ≥2 elems).
func TestNilLess(t *testing.T) {
	pq := New[int](nil)
	testutil.EqVal(t, pq.Len(), 0)
	testutil.EqVal(t, Of[int](nil).Len(), 0)
	testutil.EqVal(t, Of[int](nil, 1).Len(), 1) // single elem: Init never calls Less

	single := New[int](nil)
	single.Push(1)
	val, ok := single.Pop()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, val, 1)

	cases := []struct {
		name string
		fn   func()
	}{
		{"Of heapify", func() { Of[int](nil, 1, 2) }},
		{"FromSlice", func() { FromSlice[int](nil, []int{1, 2}) }},
		{"Push", func() {
			pq := New[int](nil)
			pq.Push(1)
			pq.Push(2)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func TestNilReceiverPanics(t *testing.T) {
	var pq *PriorityQueue[int]
	cases := []struct {
		name string
		fn   func()
	}{
		{"Push", func() { pq.Push(1) }},
		{"Pop", func() { pq.Pop() }},
		{"Peek", func() { pq.Peek() }},
		{"Remove", func() { pq.Remove(1) }},
		{"Clear", func() { pq.Clear() }},
		{"Len", func() { pq.Len() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

// eqElems checks membership (sorted by less). Does not assert heap layout.
func eqElems[T comparable](t *testing.T, pq *PriorityQueue[T], want []T) {
	t.Helper()
	got := sortedBy(pq.h.less, pq.ToSlice())
	testutil.Eq(t, got, sortedBy(pq.h.less, want))
	if len(want) == 0 {
		testutil.NonNilEmpty(t, pq.ToSlice())
	}
}

func eqPoll[T comparable](t *testing.T, pq *PriorityQueue[T], want []T) {
	t.Helper()
	testutil.Eq(t, drain(pq), want)
}

func drain[T comparable](pq *PriorityQueue[T]) []T {
	var out []T
	for {
		v, ok := pq.Pop()
		if !ok {
			break
		}
		out = append(out, v)
	}
	if out == nil {
		return []T{}
	}
	return out
}

func sortedBy[T comparable](less func(a, b T) bool, vals []T) []T {
	out := slices.Clone(vals)
	if out == nil {
		out = []T{}
	}
	slices.SortFunc(out, func(a, b T) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})
	return out
}
