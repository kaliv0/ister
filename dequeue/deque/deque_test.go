package deque

import (
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func TestNew(t *testing.T) {
	eqDeque(t, New[int](), []int{})
	testutil.EqVal(t, cap(New[int]().buf), 0)
}

func TestTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		d := Of("a", "b")
		d.PushBack("c")
		front, ok := d.PeekFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, front, "a")
		back, ok := d.PeekBack()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, back, "c")
		val, ok := d.PopFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, "a")
		eqDeque(t, d, []string{"b", "c"})
	})
	t.Run("struct", func(t *testing.T) {
		d := New[testutil.Pair]()
		d.PushBack(testutil.Pair{A: 1, B: 2})
		d.PushFront(testutil.Pair{A: 3, B: 4})
		val, ok := d.PopFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, testutil.Pair{A: 3, B: 4})
		eqDeque(t, d, []testutil.Pair{{A: 1, B: 2}})
	})
}

func TestOf(t *testing.T) {
	eqDeque(t, Of(1, 2, 3), []int{1, 2, 3})
	eqDeque(t, Of[int](), []int{})
}

func TestOfCopies(t *testing.T) {
	src := make([]int, 3, 8)
	copy(src, []int{1, 2, 3})
	d := Of(src...)
	d.PushBack(4)
	testutil.EqVal(t, src[:cap(src)][3], 0)
	eqDeque(t, d, []int{1, 2, 3, 4})
}

func TestOfCapacity(t *testing.T) {
	cases := []struct {
		name    string
		n       int
		wantCap int
	}{
		{"empty", 0, 0},
		{"one", 1, 8},
		{"eight", 8, 8},
		{"nine", 9, 16},
		{"sixteen", 16, 16},
		{"seventeen", 17, 32},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vals := make([]int, tc.n)
			for i := range vals {
				vals[i] = i + 1
			}
			d := Of(vals...)
			testutil.EqVal(t, cap(d.buf), tc.wantCap)
			testutil.EqVal(t, d.head, 0)
			testutil.EqVal(t, d.count, tc.n)
			eqDeque(t, d, vals)
		})
	}
}

func TestFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		d := FromSlice(src)
		src[0] = 99
		eqDeque(t, d, []int{1, 2, 3})
	})
	t.Run("first is front", func(t *testing.T) {
		d := FromSlice([]int{1, 2, 3})
		front, ok := d.PeekFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, front, 1)
	})
	t.Run("nil", func(t *testing.T) {
		d := FromSlice([]int(nil))
		testutil.EqVal(t, d.Len(), 0)
		testutil.EqVal(t, cap(d.buf), 0)
	})
}

func TestZeroValue(t *testing.T) {
	for _, tc := range []struct {
		name string
		push func(*Deque[int], int)
	}{
		{"PushBack", (*Deque[int]).PushBack},
		{"PushFront", (*Deque[int]).PushFront},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := &Deque[int]{}
			testutil.EqVal(t, d.Len(), 0)
			tc.push(d, 1)
			eqDeque(t, d, []int{1})
			testutil.EqVal(t, cap(d.buf), 8)
		})
	}
}

func TestLen(t *testing.T) {
	testutil.EqVal(t, New[int]().Len(), 0)
	testutil.EqVal(t, Of(1, 2, 3).Len(), 3)
}

func TestIsEmpty(t *testing.T) {
	testutil.EqVal(t, New[int]().IsEmpty(), true)
	testutil.EqVal(t, Of(1).IsEmpty(), false)
}

func TestPush(t *testing.T) {
	t.Run("back", func(t *testing.T) {
		d := New[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		eqDeque(t, d, []int{1, 2, 3})
		testutil.EqVal(t, cap(d.buf), 8)
	})
	t.Run("front", func(t *testing.T) {
		d := New[int]()
		d.PushFront(1)
		d.PushFront(2)
		d.PushFront(3)
		eqDeque(t, d, []int{3, 2, 1})
		testutil.EqVal(t, cap(d.buf), 8)
	})
	t.Run("both ends", func(t *testing.T) {
		d := New[int]()
		d.PushBack(2)
		d.PushFront(1)
		d.PushBack(3)
		d.PushFront(0)
		eqDeque(t, d, []int{0, 1, 2, 3})
	})
}

func TestPop(t *testing.T) {
	type caseT struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
		want    []int
	}
	front := []caseT{
		{"three", []int{1, 2, 3}, 1, true, []int{2, 3}},
		{"one", []int{1}, 1, true, []int{}},
		{"empty", nil, 0, false, []int{}},
	}
	back := []caseT{
		{"three", []int{1, 2, 3}, 3, true, []int{1, 2}},
		{"one", []int{1}, 1, true, []int{}},
		{"empty", nil, 0, false, []int{}},
	}
	for _, tc := range front {
		t.Run("front/"+tc.name, func(t *testing.T) {
			d := Of(tc.start...)
			val, ok := d.PopFront()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			eqDeque(t, d, tc.want)
		})
	}
	for _, tc := range back {
		t.Run("back/"+tc.name, func(t *testing.T) {
			d := Of(tc.start...)
			val, ok := d.PopBack()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			eqDeque(t, d, tc.want)
		})
	}
}

func TestDrain(t *testing.T) {
	t.Run("front", func(t *testing.T) {
		d := Of(1, 2, 3)
		var got []int
		for {
			v, ok := d.PopFront()
			if !ok {
				break
			}
			got = append(got, v)
		}
		testutil.Eq(t, got, []int{1, 2, 3})
		testutil.EqVal(t, d.head, 0)
	})
	t.Run("back", func(t *testing.T) {
		d := Of(1, 2, 3)
		var got []int
		for {
			v, ok := d.PopBack()
			if !ok {
				break
			}
			got = append(got, v)
		}
		testutil.Eq(t, got, []int{3, 2, 1})
		testutil.EqVal(t, d.head, 0)
	})
}

func TestReuseAfterPop(t *testing.T) {
	t.Run("pop front then push back", func(t *testing.T) {
		d := Of(1, 2, 3)
		d.PopFront()
		d.PopFront()
		d.PushBack(4)
		eqDeque(t, d, []int{3, 4})
	})
	t.Run("pop back then push front", func(t *testing.T) {
		d := wrapFull(t)
		val, ok := d.PopBack()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, 10)
		d.PushFront(2)
		eqDeque(t, d, []int{2, 3, 4, 5, 6, 7, 8, 9})
	})
}

func TestPeek(t *testing.T) {
	type caseT struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
	}
	front := []caseT{
		{"three", []int{1, 2, 3}, 1, true},
		{"one", []int{1}, 1, true},
		{"empty", nil, 0, false},
	}
	back := []caseT{
		{"three", []int{1, 2, 3}, 3, true},
		{"one", []int{1}, 1, true},
		{"empty", nil, 0, false},
	}
	for _, tc := range front {
		t.Run("front/"+tc.name, func(t *testing.T) {
			d := Of(tc.start...)
			val, ok := d.PeekFront()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			testutil.EqVal(t, d.Len(), len(tc.start))
		})
	}
	for _, tc := range back {
		t.Run("back/"+tc.name, func(t *testing.T) {
			d := Of(tc.start...)
			val, ok := d.PeekBack()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			testutil.EqVal(t, d.Len(), len(tc.start))
		})
	}
}

func TestClear(t *testing.T) {
	t.Run("keeps capacity", func(t *testing.T) {
		d := Of(1, 2, 3, 4, 5)
		capBefore := cap(d.buf)
		d.Clear()
		testutil.EqVal(t, d.Len(), 0)
		testutil.EqVal(t, d.head, 0)
		if cap(d.buf) != capBefore {
			t.Fatalf("Clear should keep capacity: got %d, want %d", cap(d.buf), capBefore)
		}
		d.PushBack(9)
		eqDeque(t, d, []int{9})
	})
	t.Run("resets head", func(t *testing.T) {
		d := Of(1, 2, 3)
		d.PopFront()
		d.PopFront()
		d.Clear()
		testutil.EqVal(t, d.head, 0)
		d.PushBack(9)
		front, ok := d.PeekFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, front, 9)
	})
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
			t.Fatal("empty deque should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	d := Of(1, 2, 3)
	a := d.ToSlice()
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	testutil.Eq(t, d.ToSlice(), []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := Of(1, 2)
	cleared.Clear()
	drained := Of(1, 2)
	drained.PopFront()
	drained.PopFront()
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
	var d *Deque[int]
	cases := []struct {
		name string
		fn   func()
	}{
		{"PushFront", func() { d.PushFront(1) }},
		{"PushBack", func() { d.PushBack(1) }},
		{"PopFront", func() { d.PopFront() }},
		{"PopBack", func() { d.PopBack() }},
		{"Clear", func() { d.Clear() }},
		{"Len", func() { d.Len() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func TestWrap(t *testing.T) {
	d := wrapFull(t)
	testutil.EqVal(t, d.head != 0, true)
	testutil.EqVal(t, d.count, 8)
	testutil.EqVal(t, cap(d.buf), 8)
	eqDeque(t, d, []int{3, 4, 5, 6, 7, 8, 9, 10})

	t.Run("peek and pop", func(t *testing.T) {
		d := wrapFull(t)
		front, ok := d.PeekFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, front, 3)
		back, ok := d.PeekBack()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, back, 10)

		val, ok := d.PopFront()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, 3)
		val, ok = d.PopBack()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, 10)
		eqDeque(t, d, []int{4, 5, 6, 7, 8, 9})
	})
	t.Run("All", func(t *testing.T) {
		var vals []int
		for v := range wrapFull(t).All() {
			vals = append(vals, v)
		}
		testutil.Eq(t, vals, []int{3, 4, 5, 6, 7, 8, 9, 10})
	})
	t.Run("ToSlice", func(t *testing.T) {
		testutil.Eq(t, wrapFull(t).ToSlice(), []int{3, 4, 5, 6, 7, 8, 9, 10})
	})
}

func TestGrow(t *testing.T) {
	t.Run("from PushBack", func(t *testing.T) {
		d := Of(1, 2, 3, 4, 5, 6, 7, 8)
		testutil.EqVal(t, cap(d.buf), 8)
		d.PushBack(9)
		testutil.EqVal(t, d.head, 0)
		testutil.EqVal(t, cap(d.buf), 16)
		eqDeque(t, d, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	})
	t.Run("unwraps", func(t *testing.T) {
		d := wrapFull(t)
		d.PushBack(11)
		testutil.EqVal(t, d.head, 0)
		testutil.EqVal(t, cap(d.buf), 16)
		eqDeque(t, d, []int{3, 4, 5, 6, 7, 8, 9, 10, 11})
	})
	t.Run("from PushFront", func(t *testing.T) {
		d := wrapFull(t)
		d.PushFront(2)
		testutil.EqVal(t, d.head, 15)
		testutil.EqVal(t, cap(d.buf), 16)
		eqDeque(t, d, []int{2, 3, 4, 5, 6, 7, 8, 9, 10})
	})
}

func TestEmptyResetsHead(t *testing.T) {
	d := Of(1, 2, 3)
	d.PopFront()
	testutil.EqVal(t, d.head != 0, true)
	d.PopFront()
	d.PopFront()
	testutil.EqVal(t, d.count, 0)
	testutil.EqVal(t, d.head, 0)
	d.PushBack(9)
	eqDeque(t, d, []int{9})
}

func TestPopZeros(t *testing.T) {
	t.Run("front", func(t *testing.T) {
		a, b := 1, 2
		d := Of(&a, &b)
		d.PopFront()
		testutil.EqVal(t, d.buf[0], (*int)(nil))
		eqDeque(t, d, []*int{&b})
	})
	t.Run("back", func(t *testing.T) {
		a, b := 1, 2
		d := Of(&a, &b)
		d.PopBack()
		testutil.EqVal(t, d.buf[1], (*int)(nil))
		eqDeque(t, d, []*int{&a})
	})
}

func TestClearZeros(t *testing.T) {
	a, b, c := 1, 2, 3
	d := Of(&a, &b, &c)
	d.Clear()
	for i := 0; i < cap(d.buf); i++ {
		testutil.EqVal(t, d.buf[i], (*int)(nil))
	}
}

// wrapFull builds a full ring: head != 0, n == cap == 8, logical [3..10].
func wrapFull(t *testing.T) *Deque[int] {
	t.Helper()
	d := Of(1, 2, 3, 4, 5, 6, 7, 8)
	testutil.EqVal(t, cap(d.buf), 8)
	d.PopFront()
	d.PopFront()
	d.PushBack(9)
	d.PushBack(10)
	return d
}

func eqDeque[T comparable](t *testing.T, d *Deque[T], want []T) {
	t.Helper()
	testutil.Eq(t, d.ToSlice(), want)
}
