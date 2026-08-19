package set

import (
	"slices"
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func TestNew(t *testing.T) {
	eqSet(t, New[int](), []int{})
}

func TestOf(t *testing.T) {
	eqSet(t, Of(1, 2, 3), []int{1, 2, 3})
	eqSet(t, Of[int](), []int{})
	eqSet(t, Of(1, 2, 2, 1), []int{1, 2})
}

func TestFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{1, 2, 3}
		s := FromSlice(src)
		src[0] = 99
		eqSet(t, s, []int{1, 2, 3})
		testutil.Eq(t, src, []int{99, 2, 3})
	})
	t.Run("duplicates", func(t *testing.T) {
		eqSet(t, FromSlice([]int{1, 2, 2, 3, 1}), []int{1, 2, 3})
	})
	t.Run("nil", func(t *testing.T) {
		eqSet(t, FromSlice([]int(nil)), []int{})
	})
}

func TestZeroValue(t *testing.T) {
	s := &Set[int]{}
	testutil.EqVal(t, s.Len(), 0)
	testutil.EqVal(t, s.Contains(1), false)
	s.Add(1)
	eqSet(t, s, []int{1})
}

func TestLen(t *testing.T) {
	testutil.EqVal(t, New[int]().Len(), 0)
	testutil.EqVal(t, Of(1, 2, 2).Len(), 2)
}

func TestAdd(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		s := New[int]()
		s.Add(1)
		s.Add(2, 3)
		eqSet(t, s, []int{1, 2, 3})
	})
	t.Run("no values", func(t *testing.T) {
		s := Of(1, 2)
		s.Add()
		eqSet(t, s, []int{1, 2})
	})
	t.Run("duplicates", func(t *testing.T) {
		s := Of(1)
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
		{"found", Of(1, 2, 3), 2, true},
		{"missing", Of(1, 2, 3), 9, false},
		{"empty", New[int](), 1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.s.Contains(tc.val), tc.want)
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
			s := Of(tc.start...)
			testutil.EqVal(t, s.Remove(tc.val), tc.ok)
			eqSet(t, s, tc.want)
		})
	}
}

func TestClear(t *testing.T) {
	s := Of(1, 2, 3)
	s.Clear()
	eqSet(t, s, []int{})
	s.Add(4)
	eqSet(t, s, []int{4})
}

func TestAll(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		var vals []int
		for v := range Of(1, 2, 3).All() {
			vals = append(vals, v)
		}
		slices.Sort(vals)
		testutil.Eq(t, vals, []int{1, 2, 3})
	})
	t.Run("early break", func(t *testing.T) {
		n := 0
		for range Of(1, 2, 3).All() {
			n++
			break
		}
		testutil.EqVal(t, n, 1)
	})
	t.Run("empty", func(t *testing.T) {
		for range New[int]().All() {
			t.Fatal("empty set should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	s := Of(1, 2, 3)
	a := s.ToSlice()
	slices.Sort(a)
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	got := s.ToSlice()
	slices.Sort(got)
	testutil.Eq(t, got, []int{1, 2, 3})
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
		{"empty", New[int]().String(), "{}"},
		{"single", Of(1).String(), "{1}"},
		{"stringer", Of(testutil.Ident(1)).String(), "{id=1}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.got, tc.want)
		})
	}
}

func TestUnion(t *testing.T) {
	cases := []struct {
		name string
		a, b []int
		want []int
	}{
		{"overlap", []int{1, 2}, []int{2, 3}, []int{1, 2, 3}},
		{"disjoint", []int{1}, []int{2}, []int{1, 2}},
		{"empty other", []int{1, 2}, nil, []int{1, 2}},
		{"empty self", nil, []int{1, 2}, []int{1, 2}},
		{"both empty", nil, nil, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, b := Of(tc.a...), Of(tc.b...)
			eqSet(t, a.Union(b), tc.want)
			eqSet(t, a, tc.a)
			eqSet(t, b, tc.b)
		})
	}
}

func TestIntersect(t *testing.T) {
	cases := []struct {
		name string
		a, b []int
		want []int
	}{
		{"overlap", []int{1, 2, 3}, []int{2, 3, 4}, []int{2, 3}},
		{"one sided", []int{1, 2, 3}, []int{2}, []int{2}},
		{"other larger", []int{2}, []int{1, 2, 3}, []int{2}},
		{"disjoint", []int{1}, []int{2}, []int{}},
		{"empty other", []int{1, 2}, nil, []int{}},
		{"empty self", nil, []int{1, 2}, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, b := Of(tc.a...), Of(tc.b...)
			eqSet(t, a.Intersect(b), tc.want)
			eqSet(t, a, tc.a)
			eqSet(t, b, tc.b)
		})
	}
}

func TestDifference(t *testing.T) {
	cases := []struct {
		name string
		a, b []int
		want []int
	}{
		{"partial", []int{1, 2, 3}, []int{2}, []int{1, 3}},
		{"all removed", []int{1, 2}, []int{1, 2, 3}, []int{}},
		{"empty other", []int{1, 2}, nil, []int{1, 2}},
		{"empty self", nil, []int{1}, []int{}},
		{"disjoint", []int{1, 2}, []int{3}, []int{1, 2}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, b := Of(tc.a...), Of(tc.b...)
			eqSet(t, a.Difference(b), tc.want)
			eqSet(t, a, tc.a)
			eqSet(t, b, tc.b)
		})
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		name string
		a, b *Set[int]
		want bool
	}{
		{"same", Of(1, 2), Of(2, 1), true},
		{"different", Of(1, 2), Of(1, 2, 3), false},
		{"empty", New[int](), Of[int](), true},
		{"zero and New", &Set[int]{}, New[int](), true},
		{"empty vs value", New[int](), Of(1), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.a.Equal(tc.b), tc.want)
		})
	}
}

func TestIsSubsetOf(t *testing.T) {
	cases := []struct {
		name string
		a, b *Set[int]
		want bool
	}{
		{"subset", Of(1, 2), Of(1, 2, 3), true},
		{"not subset", Of(1, 2, 3), Of(1, 2), false},
		{"equal", Of(1, 2), Of(2, 1), true},
		{"empty of values", New[int](), Of(1), true},
		{"empty of empty", New[int](), New[int](), true},
		{"value of empty", Of(1), New[int](), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.a.IsSubsetOf(tc.b), tc.want)
		})
	}
}

func TestIsSupersetOf(t *testing.T) {
	cases := []struct {
		name string
		a, b *Set[int]
		want bool
	}{
		{"superset", Of(1, 2, 3), Of(1, 2), true},
		{"not superset", Of(1, 2), Of(1, 2, 3), false},
		{"equal", Of(1, 2), Of(2, 1), true},
		{"values of empty", Of(1), New[int](), true},
		{"empty of empty", New[int](), New[int](), true},
		{"empty of values", New[int](), Of(1), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.a.IsSupersetOf(tc.b), tc.want)
		})
	}
}

func TestIsDisjoint(t *testing.T) {
	cases := []struct {
		name string
		a, b *Set[int]
		want bool
	}{
		{"disjoint", Of(1, 2), Of(3, 4), true},
		{"overlap", Of(1, 2), Of(2, 3), false},
		{"equal", Of(1, 2), Of(1, 2), false},
		{"empty other", Of(1), New[int](), true},
		{"empty self", New[int](), Of(1), true},
		{"both empty", New[int](), New[int](), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.a.IsDisjoint(tc.b), tc.want)
		})
	}
}

func TestSymmetricDifference(t *testing.T) {
	cases := []struct {
		name string
		a, b []int
		want []int
	}{
		{"overlap", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 4}},
		{"disjoint", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"equal", []int{1, 2}, []int{1, 2}, []int{}},
		{"empty other", []int{1, 2}, nil, []int{1, 2}},
		{"empty self", nil, []int{1, 2}, []int{1, 2}},
		{"both empty", nil, nil, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, b := Of(tc.a...), Of(tc.b...)
			eqSet(t, a.SymmetricDifference(b), tc.want)
			eqSet(t, a, tc.a)
			eqSet(t, b, tc.b)
		})
	}
}

func TestNilOperand(t *testing.T) {
	eqSet(t, Of(1, 2).Union(nil), []int{1, 2})
	eqSet(t, Of(1).Intersect(nil), []int{})
	eqSet(t, Of(1, 2).Difference(nil), []int{1, 2})
	eqSet(t, Of(1, 2).SymmetricDifference(nil), []int{1, 2})
	testutil.EqVal(t, Of(1).Equal(nil), false)
	testutil.EqVal(t, New[int]().Equal(nil), true)
	testutil.EqVal(t, Of(1).IsSubsetOf(nil), false)
	testutil.EqVal(t, New[int]().IsSubsetOf(nil), true)
	testutil.EqVal(t, Of(1).IsSupersetOf(nil), true)
	testutil.EqVal(t, New[int]().IsSupersetOf(nil), true)
	testutil.EqVal(t, Of(1).IsDisjoint(nil), true)
}

func TestNilReceiverPanics(t *testing.T) {
	var s *Set[int]
	cases := []struct {
		name string
		fn   func()
	}{
		{"Add", func() { s.Add(1) }},
		{"Union", func() { s.Union(Of(1)) }},
		{"Equal", func() { s.Equal(New[int]()) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
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
