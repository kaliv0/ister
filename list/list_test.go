package list

import (
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func TestNew(t *testing.T) {
	eqList(t, New[int](), []int{})
}

func TestOf(t *testing.T) {
	eqList(t, Of(1, 2, 3), []int{1, 2, 3})
	eqList(t, Of[int](), []int{})
}

func TestOfCopies(t *testing.T) {
	src := make([]int, 3, 8)
	copy(src, []int{1, 2, 3})
	l := Of(src...)
	l.Add(4)
	testutil.EqVal(t, src[:cap(src)][3], 0)
	eqList(t, l, []int{1, 2, 3, 4})
}

func TestListFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		l := ListFromSlice(src)
		src[0] = 99
		eqList(t, l, []int{1, 2, 3})
		testutil.Eq(t, src, []int{99, 2, 3})
	})
	t.Run("list mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		l := ListFromSlice(src)
		l.Set(1, 20)
		testutil.EqVal(t, src[1], 2)
		testutil.EqVal(t, l.Get(1), 20)
	})
}

func TestAdd(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		l := New[int]()
		l.Add(1)
		l.Add(2, 3)
		eqList(t, l, []int{1, 2, 3})
	})
	t.Run("no values", func(t *testing.T) {
		l := Of(1, 2)
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
			l := Of(tc.start...)
			l.Insert(tc.i, tc.vals...)
			eqList(t, l, tc.want)
		})
	}
}

func TestGet(t *testing.T) {
	l := Of(10, 20, 30)
	testutil.EqVal(t, l.Get(0), 10)
	testutil.EqVal(t, l.Get(1), 20)
	testutil.EqVal(t, l.Get(2), 30)
}

func TestSet(t *testing.T) {
	l := Of(1, 2, 3)
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
			l := Of(tc.start...)
			testutil.EqVal(t, l.RemoveAt(tc.i), tc.wantVal)
			eqList(t, l, tc.want)
		})
	}
}

func TestClear(t *testing.T) {
	l := Of(1, 2, 3, 4, 5)
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
	l := Of(&a, &b)
	old := l.items
	l.Clear()
	testutil.Eq(t, old, []*int{nil, nil})
}

func TestAll(t *testing.T) {
	t.Run("values", func(t *testing.T) {
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
			if v == 2 {
				break
			}
		}
		testutil.Eq(t, vals, []int{1, 2})
	})
	t.Run("empty", func(t *testing.T) {
		for range New[int]().All() {
			t.Fatal("empty list should not yield")
		}
	})
	t.Run("live view", func(t *testing.T) {
		l := Of(1, 2, 3)
		var vals []int
		for v := range l.All() {
			vals = append(vals, v)
			if v == 1 {
				l.Set(2, 99)
			}
		}
		testutil.Eq(t, vals, []int{1, 2, 99})
	})
}

func TestToSlice(t *testing.T) {
	l := Of(1, 2, 3)
	a := l.ToSlice()
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	testutil.Eq(t, l.ToSlice(), []int{1, 2, 3})
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
		{"ListFromSlice nil", ListFromSlice([]int(nil)).ToSlice()},
		{"after Clear", cleared.ToSlice()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.NonNilEmpty(t, tc.got)
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
			l := Of(tc.start...)
			testutil.EqVal(t, Remove(l, tc.val), tc.ok)
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
		{"found", Of(1, 2, 3), 2, 1},
		{"missing", Of(1, 2, 3), 9, -1},
		{"empty", New[int](), 1, -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, Index(tc.l, tc.val), tc.want)
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
		{"found", Of(1, 2, 3), 2, true},
		{"missing", Of(1, 2, 3), 9, false},
		{"empty", New[int](), 1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, Contains(tc.l, tc.val), tc.want)
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
			l := Of(tc.start...)
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
			l := Of(tc.start...)
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

func TestOutOfRangePanics(t *testing.T) {
	l := Of(1)
	empty := New[int]()
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
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func eqList[T comparable](t *testing.T, l *List[T], want []T) {
	t.Helper()
	testutil.Eq(t, l.items, want)
}
