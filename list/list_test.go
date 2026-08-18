package list

import (
	"fmt"
	"slices"
	"testing"
)

func TestNewList(t *testing.T) {
	eqList(t, NewList[int](), []int{})
}

func TestOf(t *testing.T) {
	eqList(t, ListOf(1, 2, 3), []int{1, 2, 3})
	eqList(t, ListOf[int](), []int{})
}

func TestOfCopies(t *testing.T) {
	src := make([]int, 3, 8)
	copy(src, []int{1, 2, 3})
	l := ListOf(src...)
	l.Add(4)
	eqVal(t, src[:cap(src)][3], 0)
	eqList(t, l, []int{1, 2, 3, 4})
}

func TestListFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		l := ListFromSlice(src)
		src[0] = 99
		eqList(t, l, []int{1, 2, 3})
		eq(t, src, []int{99, 2, 3})
	})
	t.Run("list mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		l := ListFromSlice(src)
		l.Set(1, 20)
		eqVal(t, src[1], 2)
		eqVal(t, l.Get(1), 20)
	})
}

func TestAdd(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		l := NewList[int]()
		l.Add(1)
		l.Add(2, 3)
		eqList(t, l, []int{1, 2, 3})
	})
	t.Run("no values", func(t *testing.T) {
		l := ListOf(1, 2)
		l.Add()
		eqList(t, l, []int{1, 2})
	})
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name  string
		start []int
		i     int
		vals  []int
		want  []int
	}{
		{"middle", []int{1, 3}, 1, []int{2}, []int{1, 2, 3}},
		{"append", []int{1, 2, 3}, 3, []int{4}, []int{1, 2, 3, 4}},
		{"multiple", []int{1, 2, 3, 4}, 2, []int{20, 30, 40}, []int{1, 2, 20, 30, 40, 3, 4}},
		{"front", []int{1, 2, 3}, 0, []int{0}, []int{0, 1, 2, 3}},
		{"no values", []int{1, 2, 3}, 1, nil, []int{1, 2, 3}},
		{"into empty", nil, 0, []int{7}, []int{7}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := ListOf(tc.start...)
			l.Insert(tc.i, tc.vals...)
			eqList(t, l, tc.want)
		})
	}
}

func TestGet(t *testing.T) {
	l := ListOf(10, 20, 30)
	eqVal(t, l.Get(0), 10)
	eqVal(t, l.Get(1), 20)
	eqVal(t, l.Get(2), 30)
}

func TestSet(t *testing.T) {
	l := ListOf(1, 2, 3)
	l.Set(0, 10)
	l.Set(1, 20)
	l.Set(2, 30)
	eqList(t, l, []int{10, 20, 30})
}

func TestRemoveAt(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		i       int
		wantVal int
		want    []int
	}{
		{"middle", []int{1, 2, 3}, 1, 2, []int{1, 3}},
		{"first", []int{1, 2, 3}, 0, 1, []int{2, 3}},
		{"last", []int{1, 2, 3}, 2, 3, []int{1, 2}},
		{"only", []int{1}, 0, 1, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := ListOf(tc.start...)
			eqVal(t, l.RemoveAt(tc.i), tc.wantVal)
			eqList(t, l, tc.want)
		})
	}
}

func TestClear(t *testing.T) {
	l := ListOf(1, 2, 3, 4, 5)
	capBefore := cap(l.items)
	l.Clear()
	eqList(t, l, []int{})
	if cap(l.items) != capBefore {
		t.Fatalf("Clear should keep capacity: got %d, want %d", cap(l.items), capBefore)
	}

	l.Add(3)
	eqList(t, l, []int{3})
}

func TestClearZeros(t *testing.T) {
	a, b := 1, 2
	l := ListOf(&a, &b)
	old := l.items
	l.Clear()
	eq(t, old, []*int{nil, nil})
}

func TestAll(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		var vals []int
		for v := range ListOf(1, 2, 3).All() {
			vals = append(vals, v)
		}
		eq(t, vals, []int{1, 2, 3})
	})
	t.Run("early break", func(t *testing.T) {
		var vals []int
		for v := range ListOf(1, 2, 3).All() {
			vals = append(vals, v)
			if v == 2 {
				break
			}
		}
		eq(t, vals, []int{1, 2})
	})
	t.Run("empty", func(t *testing.T) {
		for range NewList[int]().All() {
			t.Fatal("empty list should not yield")
		}
	})
	t.Run("live view", func(t *testing.T) {
		l := ListOf(1, 2, 3)
		var vals []int
		for v := range l.All() {
			vals = append(vals, v)
			if v == 1 {
				l.Set(2, 99)
			}
		}
		eq(t, vals, []int{1, 2, 99})
	})
}

func TestToSlice(t *testing.T) {
	l := ListOf(1, 2, 3)
	a := l.ToSlice()
	eq(t, a, []int{1, 2, 3})

	a[0] = 99
	eq(t, l.ToSlice(), []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := ListOf(1, 2)
	cleared.Clear()
	cases := []struct {
		name string
		got  []int
	}{
		{"NewList", NewList[int]().ToSlice()},
		{"ListOf", ListOf[int]().ToSlice()},
		{"ListFromSlice nil", ListFromSlice([]int(nil)).ToSlice()},
		{"after Clear", cleared.ToSlice()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nonNilEmpty(t, tc.got)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name  string
		start []string
		val   string
		ok    bool
		want  []string
	}{
		{"first match", []string{"a", "b", "b"}, "b", true, []string{"a", "b"}},
		{"missing", []string{"a", "b"}, "z", false, []string{"a", "b"}},
		{"first", []string{"a", "b", "c"}, "a", true, []string{"b", "c"}},
		{"last", []string{"a", "b", "c"}, "c", true, []string{"a", "b"}},
		{"only", []string{"x"}, "x", true, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := ListOf(tc.start...)
			eqVal(t, Remove(l, tc.val), tc.ok)
			eqList(t, l, tc.want)
		})
	}
}

func TestIndex(t *testing.T) {
	cases := []struct {
		name string
		l    *List[int]
		val  int
		want int
	}{
		{"found", ListOf(1, 2, 3), 2, 1},
		{"missing", ListOf(1, 2, 3), 9, -1},
		{"empty", NewList[int](), 1, -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eqVal(t, Index(tc.l, tc.val), tc.want)
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		name string
		l    *List[int]
		val  int
		want bool
	}{
		{"found", ListOf(1, 2, 3), 2, true},
		{"missing", ListOf(1, 2, 3), 9, false},
		{"empty", NewList[int](), 1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eqVal(t, Contains(tc.l, tc.val), tc.want)
		})
	}
}

func TestReverse(t *testing.T) {
	cases := []struct {
		name  string
		start []int
		want  []int
	}{
		{"three", []int{1, 2, 3}, []int{3, 2, 1}},
		{"empty", nil, []int{}},
		{"single", []int{1}, []int{1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := ListOf(tc.start...)
			l.Reverse()
			eqList(t, l, tc.want)
		})
	}
}

func TestSort(t *testing.T) {
	cases := []struct {
		name  string
		start []int
		want  []int
	}{
		{"unsorted", []int{3, 1, 2}, []int{1, 2, 3}},
		{"empty", nil, []int{}},
		{"already sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"duplicates", []int{2, 1, 2, 1}, []int{1, 1, 2, 2}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := ListOf(tc.start...)
			Sort(l)
			eqList(t, l, tc.want)
		})
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ints", ListOf(1, 2, 3).String(), "[1, 2, 3]"},
		{"empty", NewList[int]().String(), "[]"},
		{"stringer", ListOf(ident(1), ident(2)).String(), "[id=1, id=2]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eqVal(t, tc.got, tc.want)
		})
	}
}

func TestOutOfRangePanics(t *testing.T) {
	l := ListOf(1)
	empty := NewList[int]()
	cases := []struct {
		name string
		fn   func()
	}{
		{"Get negative", func() { l.Get(-1) }},
		{"Get past end", func() { l.Get(10) }},
		{"Set negative", func() { l.Set(-1, 0) }},
		{"Set past end", func() { l.Set(10, 0) }},
		{"Insert negative", func() { l.Insert(-1, 0) }},
		{"Insert past end", func() { l.Insert(10, 0) }},
		{"RemoveAt negative", func() { l.RemoveAt(-1) }},
		{"RemoveAt past end", func() { l.RemoveAt(10) }},
		{"Get empty", func() { empty.Get(0) }},
		{"Set empty", func() { empty.Set(0, 1) }},
		{"RemoveAt empty", func() { empty.RemoveAt(0) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mustPanic(t, tc.fn)
		})
	}
}

func eq[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func eqList[T comparable](t *testing.T, l *List[T], want []T) {
	t.Helper()
	eq(t, l.items, want)
}

func eqVal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func nonNilEmpty[T any](t *testing.T, got []T) {
	t.Helper()
	if got == nil || len(got) != 0 {
		t.Fatalf("got %#v, want non-nil empty slice", got)
	}
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

type ident int

func (id ident) String() string {
	return fmt.Sprintf("id=%d", id)
}
