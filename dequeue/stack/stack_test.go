package stack

import (
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func TestStack(t *testing.T) {
	s := New[int]()
	testutil.EqVal(t, s.Len(), 0)
	testutil.EqVal(t, s.IsEmpty(), true)
}

func TestTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		s := Of("a", "b")
		s.Push("c")
		top, ok := s.Peek()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, top, "c")
		val, ok := s.Pop()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, "c")
		eqStack(t, s, []string{"a", "b"})
	})
	t.Run("struct", func(t *testing.T) {
		s := New[testutil.Pair]()
		s.Push(testutil.Pair{A: 1, B: 2})
		s.Push(testutil.Pair{A: 3, B: 4})
		val, ok := s.Pop()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, val, testutil.Pair{A: 3, B: 4})
		eqStack(t, s, []testutil.Pair{{A: 1, B: 2}})
	})
}

func TestOf(t *testing.T) {
	eqStack(t, Of(1, 2, 3), []int{1, 2, 3})
	eqStack(t, Of[int](), []int{})
}

func TestOfCopies(t *testing.T) {
	src := make([]int, 3, 8)
	copy(src, []int{1, 2, 3})
	s := Of(src...)
	s.Push(4)
	testutil.EqVal(t, src[:cap(src)][3], 0)
	eqStack(t, s, []int{1, 2, 3, 4})
}

func TestOfOrder(t *testing.T) {
	s := Of(1, 2, 3)
	top, ok := s.Peek()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, top, 3)
}

func TestFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		s := FromSlice(src)
		src[0] = 99
		eqStack(t, s, []int{1, 2, 3})
	})
	t.Run("last is top", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		top, ok := s.Peek()
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, top, 3)
	})
	t.Run("nil", func(t *testing.T) {
		s := FromSlice([]int(nil))
		testutil.EqVal(t, s.Len(), 0)
	})
}

func TestZeroValue(t *testing.T) {
	s := &Stack[int]{}
	testutil.EqVal(t, s.Len(), 0)
	testutil.EqVal(t, s.IsEmpty(), true)
	s.Push(1)
	eqStack(t, s, []int{1})
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
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	eqStack(t, s, []int{1, 2, 3})
}

func TestPop(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
		want    []int
	}{
		{"three", []int{1, 2, 3}, 3, true, []int{1, 2}},
		{"one", []int{1}, 1, true, []int{}},
		{"empty", nil, 0, false, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(tc.start...)
			val, ok := s.Pop()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			eqStack(t, s, tc.want)
		})
	}
}

func TestPopOrder(t *testing.T) {
	s := Of(1, 2, 3)
	var got []int
	for {
		v, ok := s.Pop()
		if !ok {
			break
		}
		got = append(got, v)
	}
	testutil.Eq(t, got, []int{3, 2, 1})
}

func TestPeek(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		wantVal int
		wantOk  bool
	}{
		{"three", []int{1, 2, 3}, 3, true},
		{"one", []int{1}, 1, true},
		{"empty", nil, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(tc.start...)
			val, ok := s.Peek()
			testutil.EqVal(t, ok, tc.wantOk)
			testutil.EqVal(t, val, tc.wantVal)
			testutil.EqVal(t, s.Len(), len(tc.start))
		})
	}
}

func TestClear(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	capBefore := cap(s.items)
	s.Clear()
	testutil.EqVal(t, s.Len(), 0)
	testutil.EqVal(t, s.IsEmpty(), true)
	if cap(s.items) != capBefore {
		t.Fatalf("Clear should keep capacity: got %d, want %d", cap(s.items), capBefore)
	}
	s.Push(9)
	eqStack(t, s, []int{9})
}

func TestAll(t *testing.T) {
	t.Run("top to bottom", func(t *testing.T) {
		var vals []int
		for v := range Of(1, 2, 3).All() {
			vals = append(vals, v)
		}
		testutil.Eq(t, vals, []int{3, 2, 1})
	})
	t.Run("early break", func(t *testing.T) {
		var vals []int
		for v := range Of(1, 2, 3).All() {
			vals = append(vals, v)
			break
		}
		testutil.Eq(t, vals, []int{3})
	})
	t.Run("empty", func(t *testing.T) {
		for range New[int]().All() {
			t.Fatal("empty stack should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	s := Of(1, 2, 3)
	a := s.ToSlice()
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	testutil.Eq(t, s.ToSlice(), []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := Of(1, 2)
	cleared.Clear()
	cases := []struct {
		name string
		got  []int
	}{
		{"New", New[int]().ToSlice()},
		{"Of", Of[int]().ToSlice()},
		{"FromSlice nil", FromSlice([]int(nil)).ToSlice()},
		{"after Clear", cleared.ToSlice()},
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
		{"ints", Of(1, 2, 3).String(), "[3, 2, 1]"},
		{"empty", New[int]().String(), "[]"},
		{"stringer", Of(testutil.Ident(1), testutil.Ident(2)).String(), "[id=2, id=1]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.got, tc.want)
		})
	}
}

func TestNilReceiverPanics(t *testing.T) {
	var s *Stack[int]
	cases := []struct {
		name string
		fn   func()
	}{
		{"Push", func() { s.Push(1) }},
		{"Len", func() { s.Len() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func eqStack[T comparable](t *testing.T, s *Stack[T], want []T) {
	t.Helper()
	testutil.Eq(t, s.items, want)
}
