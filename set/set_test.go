package set

import (
	"fmt"
	"slices"
	"testing"
)

func TestNewSet(t *testing.T) {
	eqSet(t, NewSet[int](), []int{})
}

func TestSetOf(t *testing.T) {
	eqSet(t, SetOf(1, 2, 3), []int{1, 2, 3})
	eqSet(t, SetOf[int](), []int{})
	eqSet(t, SetOf(1, 2, 2, 1), []int{1, 2})
}

func TestSetFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		s := SetFromSlice(src)
		src[0] = 99
		eqSet(t, s, []int{1, 2, 3})
		eq(t, src, []int{99, 2, 3})
	})
	t.Run("duplicates", func(t *testing.T) {
		eqSet(t, SetFromSlice([]int{1, 2, 2, 3, 1}), []int{1, 2, 3})
	})
	t.Run("nil", func(t *testing.T) {
		eqSet(t, SetFromSlice([]int(nil)), []int{})
	})
}

func TestZeroValue(t *testing.T) {
	s := &Set[int]{}
	eqVal(t, s.Len(), 0)
	eqVal(t, s.Contains(1), false)
	s.Add(1)
	eqSet(t, s, []int{1})
}

func TestLen(t *testing.T) {
	eqVal(t, NewSet[int]().Len(), 0)
	eqVal(t, SetOf(1, 2, 2).Len(), 2)
}

func TestAdd(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		s := NewSet[int]()
		s.Add(1)
		s.Add(2, 3)
		eqSet(t, s, []int{1, 2, 3})
	})
	t.Run("no values", func(t *testing.T) {
		s := SetOf(1, 2)
		s.Add()
		eqSet(t, s, []int{1, 2})
	})
	t.Run("duplicates", func(t *testing.T) {
		s := SetOf(1)
		s.Add(1, 2, 1)
		eqSet(t, s, []int{1, 2})
	})
}

func TestContains(t *testing.T) {
	cases := []struct {
		name string
		s    *Set[int]
		val  int
		want bool
	}{
		{"found", SetOf(1, 2, 3), 2, true},
		{"missing", SetOf(1, 2, 3), 9, false},
		{"empty", NewSet[int](), 1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eqVal(t, tc.s.Contains(tc.val), tc.want)
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
		{"present", []int{1, 2, 3}, 2, true, []int{1, 3}},
		{"missing", []int{1, 2}, 9, false, []int{1, 2}},
		{"only", []int{1}, 1, true, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := SetOf(tc.start...)
			eqVal(t, s.Remove(tc.val), tc.ok)
			eqSet(t, s, tc.want)
		})
	}
}

func TestClear(t *testing.T) {
	s := SetOf(1, 2, 3)
	s.Clear()
	eqSet(t, s, []int{})
	s.Add(4)
	eqSet(t, s, []int{4})
}

func TestAll(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		var vals []int
		for v := range SetOf(1, 2, 3).All() {
			vals = append(vals, v)
		}
		slices.Sort(vals)
		eq(t, vals, []int{1, 2, 3})
	})
	t.Run("early break", func(t *testing.T) {
		n := 0
		for range SetOf(1, 2, 3).All() {
			n++
			break
		}
		eqVal(t, n, 1)
	})
	t.Run("empty", func(t *testing.T) {
		for range NewSet[int]().All() {
			t.Fatal("empty set should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	s := SetOf(1, 2, 3)
	a := s.ToSlice()
	slices.Sort(a)
	eq(t, a, []int{1, 2, 3})

	a[0] = 99
	got := s.ToSlice()
	slices.Sort(got)
	eq(t, got, []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := SetOf(1, 2)
	cleared.Clear()
	cases := []struct {
		name string
		got  []int
	}{
		{"NewSet", NewSet[int]().ToSlice()},
		{"SetOf", SetOf[int]().ToSlice()},
		{"SetFromSlice nil", SetFromSlice([]int(nil)).ToSlice()},
		{"after Clear", cleared.ToSlice()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nonNilEmpty(t, tc.got)
		})
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"empty", NewSet[int]().String(), "{}"},
		{"single", SetOf(1).String(), "{1}"},
		{"stringer", SetOf(ident(1)).String(), "{id=1}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eqVal(t, tc.got, tc.want)
		})
	}
}

// func TestUnion(t *testing.T) {
// 	cases := []struct {
// 		name string
// 		a, b []int
// 		want []int
// 	}{
// 		{"overlap", []int{1, 2}, []int{2, 3}, []int{1, 2, 3}},
// 		{"disjoint", []int{1}, []int{2}, []int{1, 2}},
// 		{"empty other", []int{1, 2}, nil, []int{1, 2}},
// 		{"empty self", nil, []int{1, 2}, []int{1, 2}},
// 		{"both empty", nil, nil, []int{}},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			a, b := SetOf(tc.a...), SetOf(tc.b...)
// 			eqSet(t, a.Union(b), tc.want)
// 			eqSet(t, a, tc.a)
// 			eqSet(t, b, tc.b)
// 		})
// 	}
// }

// func TestIntersect(t *testing.T) {
// 	cases := []struct {
// 		name string
// 		a, b []int
// 		want []int
// 	}{
// 		{"overlap", []int{1, 2, 3}, []int{2, 3, 4}, []int{2, 3}},
// 		{"one sided", []int{1, 2, 3}, []int{2}, []int{2}},
// 		{"other larger", []int{2}, []int{1, 2, 3}, []int{2}},
// 		{"disjoint", []int{1}, []int{2}, []int{}},
// 		{"empty other", []int{1, 2}, nil, []int{}},
// 		{"empty self", nil, []int{1, 2}, []int{}},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			a, b := SetOf(tc.a...), SetOf(tc.b...)
// 			eqSet(t, a.Intersect(b), tc.want)
// 			eqSet(t, a, tc.a)
// 			eqSet(t, b, tc.b)
// 		})
// 	}
// }

// func TestDifference(t *testing.T) {
// 	cases := []struct {
// 		name string
// 		a, b []int
// 		want []int
// 	}{
// 		{"partial", []int{1, 2, 3}, []int{2}, []int{1, 3}},
// 		{"all removed", []int{1, 2}, []int{1, 2, 3}, []int{}},
// 		{"empty other", []int{1, 2}, nil, []int{1, 2}},
// 		{"empty self", nil, []int{1}, []int{}},
// 		{"disjoint", []int{1, 2}, []int{3}, []int{1, 2}},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			a, b := SetOf(tc.a...), SetOf(tc.b...)
// 			eqSet(t, a.Difference(b), tc.want)
// 			eqSet(t, a, tc.a)
// 			eqSet(t, b, tc.b)
// 		})
// 	}
// }

// func TestEqual(t *testing.T) {
// 	cases := []struct {
// 		name string
// 		a, b *Set[int]
// 		want bool
// 	}{
// 		{"same", SetOf(1, 2), SetOf(2, 1), true},
// 		{"different", SetOf(1, 2), SetOf(1, 2, 3), false},
// 		{"empty", NewSet[int](), SetOf[int](), true},
// 		{"zero and NewSet", &Set[int]{}, NewSet[int](), true},
// 		{"empty vs value", NewSet[int](), SetOf(1), false},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			eqVal(t, tc.a.Equal(tc.b), tc.want)
// 		})
// 	}
// }

// func TestIsSubsetOf(t *testing.T) {
// 	cases := []struct {
// 		name string
// 		a, b *Set[int]
// 		want bool
// 	}{
// 		{"subset", SetOf(1, 2), SetOf(1, 2, 3), true},
// 		{"not subset", SetOf(1, 2, 3), SetOf(1, 2), false},
// 		{"equal", SetOf(1, 2), SetOf(2, 1), true},
// 		{"empty of values", NewSet[int](), SetOf(1), true},
// 		{"empty of empty", NewSet[int](), NewSet[int](), true},
// 		{"value of empty", SetOf(1), NewSet[int](), false},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			eqVal(t, tc.a.IsSubsetOf(tc.b), tc.want)
// 		})
// 	}
// }

func TestNilPanics(t *testing.T) {
	var s *Set[int]
	// other := SetOf(1)
	cases := []struct {
		name string
		fn   func()
	}{
		{"Add", func() { s.Add(1) }},
		// {"Union nil other", func() { other.Union(nil) }},
		// {"Equal nil other", func() { other.Equal(nil) }},
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

func eqSet[T comparable](t *testing.T, s *Set[T], want []T) {
	t.Helper()
	if len(s.data) != len(want) {
		t.Fatalf("got %v, want %v", s.data, want)
	}
	for _, v := range want {
		if _, ok := s.data[v]; !ok {
			t.Fatalf("got %v, want %v", s.data, want)
		}
	}
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
