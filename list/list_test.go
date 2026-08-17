package list

import (
	"fmt"
	"slices"
	"testing"
)

func eq[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNew(t *testing.T) {
	if New[int]().Len() != 0 {
		t.Fatal("new list should be empty")
	}
}

func TestOf(t *testing.T) {
	eq(t, Of(1, 2, 3).ToSlice(), []int{1, 2, 3})
	if Of[int]().Len() != 0 {
		t.Fatal("Of() with no args should be empty")
	}
}

func TestFromSlice(t *testing.T) {
	src := []int{1, 2, 3}
	l := FromSlice(src)
	eq(t, l.ToSlice(), []int{1, 2, 3})

	// mutation shouldn't change original
	src[0] = 99
	eq(t, l.ToSlice(), []int{1, 2, 3})
}

func TestFromSliceEmpty(t *testing.T) {
	got := FromSlice([]int(nil)).ToSlice()
	if got == nil || len(got) != 0 {
		t.Fatalf("got %#v, want non-nil empty slice", got)
	}
}

func TestAdd(t *testing.T) {
	l := New[int]()
	l.Add(1)
	l.Add(2, 3)
	eq(t, l.ToSlice(), []int{1, 2, 3})
	if l.Len() == 0 {
		t.Fatal("Len should be non-zero after Add")
	}
}

func TestInsert(t *testing.T) {
	l := Of(1, 3)
	l.Insert(1, 2)
	eq(t, l.ToSlice(), []int{1, 2, 3})

	l.Insert(l.Len(), 4)
	eq(t, l.ToSlice(), []int{1, 2, 3, 4})

	l.Insert(2, 20, 30, 40)
	eq(t, l.ToSlice(), []int{1, 2, 20, 30, 40, 3, 4})
}

func TestGet(t *testing.T) {
	if Of(10, 20, 30).Get(1) != 20 {
		t.Fatal("Get(1) should be 20")
	}
}

func TestSet(t *testing.T) {
	l := Of(1, 2, 3)
	l.Set(1, 20)
	eq(t, l.ToSlice(), []int{1, 20, 3})
}

func TestRemoveAt(t *testing.T) {
	l := Of(1, 2, 3)
	if l.RemoveAt(1) != 2 {
		t.Fatal("RemoveAt should return 2")
	}
	eq(t, l.ToSlice(), []int{1, 3})
}

func TestClear(t *testing.T) {
	l := Of(1, 2)
	l.Clear()
	if l.Len() != 0 {
		t.Fatal("expected empty")
	}

	l.Add(3)
	eq(t, l.ToSlice(), []int{3})
}

func TestAll(t *testing.T) {
	var idxs, vals []int
	for i, v := range Of(1, 2, 3).All() {
		idxs = append(idxs, i)
		vals = append(vals, v)
	}
	eq(t, idxs, []int{0, 1, 2})
	eq(t, vals, []int{1, 2, 3})

	idxs, vals = nil, nil
	for i, v := range Of(1, 2, 3).All() {
		idxs = append(idxs, i)
		vals = append(vals, v)
		if v == 2 {
			break
		}
	}
	eq(t, idxs, []int{0, 1})
	eq(t, vals, []int{1, 2})
}

func TestToSlice(t *testing.T) {
	l := Of(1, 2, 3)
	a := l.ToSlice()
	eq(t, a, []int{1, 2, 3})

	// mutating doesn't affect the original
	a[0] = 99
	eq(t, a, []int{99, 2, 3})
	eq(t, l.ToSlice(), []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	for _, got := range [][]int{New[int]().ToSlice(), Of[int]().ToSlice()} {
		if got == nil || len(got) != 0 {
			t.Fatalf("got %#v, want non-nil empty slice", got)
		}
	}
}

func TestRemove(t *testing.T) {
	l := Of("a", "b", "b")
	if !Remove(l, "b") {
		t.Fatal("expected a match")
	}
	eq(t, l.ToSlice(), []string{"a", "b"})

	if Remove(l, "z") {
		t.Fatal("missing value should return false")
	}
	eq(t, l.ToSlice(), []string{"a", "b"})
}

func TestIndex(t *testing.T) {
	l := Of(1, 2, 3)
	if Index(l, 2) != 1 || Index(l, 9) != -1 {
		t.Fatal("Index: want 1 for 2, -1 for missing")
	}
}

func TestContains(t *testing.T) {
	l := Of(1, 2, 3)
	if !Contains(l, 2) || Contains(l, 9) {
		t.Fatal("Contains: want true for 2, false for 9")
	}
}

func TestReverse(t *testing.T) {
	l := Of(1, 2, 3)
	l.Reverse()
	eq(t, l.ToSlice(), []int{3, 2, 1})

	l2 := Of("a", "b", "c")
	l2.Reverse()
	eq(t, l2.ToSlice(), []string{"c", "b", "a"})
}

func TestSort(t *testing.T) {
	l := Of(3, 1, 2)
	Sort(l)
	eq(t, l.ToSlice(), []int{1, 2, 3})
}

func TestAddEmpty(t *testing.T) {
	l := Of(1, 2)
	l.Add()
	eq(t, l.ToSlice(), []int{1, 2})
}

func TestString(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ints", Of(1, 2, 3).String(), "[1, 2, 3]"},
		{"empty New", New[int]().String(), "[]"},
		{"empty Of", Of[int]().String(), "[]"},
		{"stringer", Of(ident(1), ident(2)).String(), "[id=1, id=2]"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("%s: got %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

func TestOutOfRangePanics(t *testing.T) {
	l := Of(1)
	mustPanic(t, func() { l.Get(-1) })
	mustPanic(t, func() { l.Get(10) })
	mustPanic(t, func() { l.Set(-1, 0) })
	mustPanic(t, func() { l.Set(10, 0) })
	mustPanic(t, func() { l.Insert(-1, 0) })
	mustPanic(t, func() { l.Insert(10, 0) })
	mustPanic(t, func() { l.RemoveAt(-1) })
	mustPanic(t, func() { l.RemoveAt(10) })
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

// ident helper
type ident int

func (id ident) String() string {
	return fmt.Sprintf("id=%d", id)
}
