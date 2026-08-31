package ordered

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func TestNew(t *testing.T) {
	eqSet(t, New(testutil.LessInt), []int{})
}

func TestTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		less := func(a, b string) bool { return a < b }
		s := Of(less, "b", "a", "b")
		eqSet(t, s, []string{"a", "b"})
		testutil.EqVal(t, s.Contains("b"), true)
		testutil.EqVal(t, s.Remove("a"), true)
		eqSet(t, s, []string{"b"})
	})
	t.Run("struct", func(t *testing.T) {
		s := Of(testutil.LessPair,
			testutil.Pair{A: 3, B: 4},
			testutil.Pair{A: 1, B: 2},
			testutil.Pair{A: 1, B: 2},
		)
		eqSet(t, s, []testutil.Pair{{A: 1, B: 2}, {A: 3, B: 4}})
		eqSet(t, s.Union(Of(testutil.LessPair, testutil.Pair{A: 3, B: 4}, testutil.Pair{A: 5, B: 6})),
			[]testutil.Pair{{A: 1, B: 2}, {A: 3, B: 4}, {A: 5, B: 6}})
	})
}

func TestOf(t *testing.T) {
	eqSet(t, Of(testutil.LessInt, 3, 1, 2), []int{1, 2, 3})
	eqSet(t, Of(testutil.LessInt), []int{})
	eqSet(t, Of(testutil.LessInt, 1, 2, 2, 1), []int{1, 2})
}

func TestFromSlice(t *testing.T) {
	t.Run("src mutation", func(t *testing.T) {
		src := []int{3, 1, 2}
		s := FromSlice(testutil.LessInt, src)
		src[0] = 99
		eqSet(t, s, []int{1, 2, 3})
		testutil.Eq(t, src, []int{99, 1, 2})
	})
	t.Run("duplicates", func(t *testing.T) {
		eqSet(t, FromSlice(testutil.LessInt, []int{1, 2, 2, 3, 1}), []int{1, 2, 3})
	})
	t.Run("nil", func(t *testing.T) {
		eqSet(t, FromSlice(testutil.LessInt, []int(nil)), []int{})
	})
}

func TestLen(t *testing.T) {
	testutil.EqVal(t, New(testutil.LessInt).Len(), 0)
	testutil.EqVal(t, Of(testutil.LessInt, 1, 2, 2).Len(), 2)
}

func TestAdd(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		s := New(testutil.LessInt)
		s.Add(3)
		s.Add(1, 2)
		eqSet(t, s, []int{1, 2, 3})
	})
	t.Run("no values", func(t *testing.T) {
		s := Of(testutil.LessInt, 1, 2)
		s.Add()
		eqSet(t, s, []int{1, 2})
	})
	t.Run("duplicates", func(t *testing.T) {
		s := Of(testutil.LessInt, 1)
		s.Add(1, 2, 1)
		eqSet(t, s, []int{1, 2})
	})
}

func TestOrder(t *testing.T) {
	t.Run("ascending", func(t *testing.T) {
		eqSet(t, Of(testutil.LessInt, 5, 1, 4, 2, 3), []int{1, 2, 3, 4, 5})
	})
	t.Run("descending", func(t *testing.T) {
		eqSet(t, Of(testutil.GreaterInt, 5, 1, 4, 2, 3), []int{5, 4, 3, 2, 1})
	})
}

func TestLessEquality(t *testing.T) {
	// Uniqueness follows less: same A means same key.
	lessA := func(a, b testutil.Pair) bool { return a.A < b.A }
	s := New(lessA)
	s.Add(testutil.Pair{A: 1, B: 10})
	s.Add(testutil.Pair{A: 1, B: 20})
	testutil.EqVal(t, s.Len(), 1)
	testutil.EqVal(t, s.Contains(testutil.Pair{A: 1, B: 99}), true)
	testutil.EqVal(t, s.Contains(testutil.Pair{A: 2, B: 1}), false)
}

func TestContains(t *testing.T) {
	cases := []struct {
		name string
		s    *OrderedSet[int]
		val  int
		want bool
	}{
		{"found", Of(testutil.LessInt, 1, 2, 3), 2, true},
		{"missing", Of(testutil.LessInt, 1, 2, 3), 9, false},
		{"empty", New(testutil.LessInt), 1, false},
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
		{"empty", nil, 1, false, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(testutil.LessInt, tc.start...)
			testutil.EqVal(t, s.Remove(tc.val), tc.ok)
			eqSet(t, s, tc.want)
		})
	}
}

func TestClear(t *testing.T) {
	s := Of(testutil.LessInt, 1, 2, 3)
	s.Clear()
	eqSet(t, s, []int{})
	s.Add(4)
	s.Add(2)
	eqSet(t, s, []int{2, 4})
}

func TestAll(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		var vals []int
		for v := range Of(testutil.LessInt, 3, 1, 2).All() {
			vals = append(vals, v)
		}
		testutil.Eq(t, vals, []int{1, 2, 3})
	})
	t.Run("early break", func(t *testing.T) {
		var vals []int
		for v := range Of(testutil.LessInt, 3, 1, 2).All() {
			vals = append(vals, v)
			break
		}
		testutil.Eq(t, vals, []int{1})
	})
	t.Run("empty", func(t *testing.T) {
		for range New(testutil.LessInt).All() {
			t.Fatal("empty set should not yield")
		}
	})
}

func TestToSlice(t *testing.T) {
	s := Of(testutil.LessInt, 3, 1, 2)
	a := s.ToSlice()
	testutil.Eq(t, a, []int{1, 2, 3})

	a[0] = 99
	testutil.Eq(t, s.ToSlice(), []int{1, 2, 3})
}

func TestToSliceEmpty(t *testing.T) {
	cleared := Of(testutil.LessInt, 1, 2)
	cleared.Clear()
	cases := []struct {
		name string
		got  []int
	}{
		{"New", New(testutil.LessInt).ToSlice()},
		{"Of", Of(testutil.LessInt).ToSlice()},
		{"FromSlice nil", FromSlice(testutil.LessInt, []int(nil)).ToSlice()},
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
		{"empty", New(testutil.LessInt).String(), "{}"},
		{"single", Of(testutil.LessInt, 1).String(), "{1}"},
		{"ascending", Of(testutil.LessInt, 3, 1, 2).String(), "{1, 2, 3}"},
		{"stringer", Of(func(a, b testutil.Ident) bool { return a < b }, testutil.Ident(1)).String(), "{id=1}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.got, tc.want)
		})
	}
}

func TestFirstLast(t *testing.T) {
	cases := []struct {
		name    string
		start   []int
		first   int
		firstOk bool
		last    int
		lastOk  bool
	}{
		{"three", []int{3, 1, 2}, 1, true, 3, true},
		{"one", []int{7}, 7, true, 7, true},
		{"empty", nil, 0, false, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(testutil.LessInt, tc.start...)
			first, ok := s.First()
			testutil.EqVal(t, ok, tc.firstOk)
			testutil.EqVal(t, first, tc.first)
			last, ok := s.Last()
			testutil.EqVal(t, ok, tc.lastOk)
			testutil.EqVal(t, last, tc.last)
		})
	}
}

func TestNavigable(t *testing.T) {
	s := Of(testutil.LessInt, 10, 20, 30, 40, 50)

	cases := []struct {
		name string
		fn   func() (int, bool)
		val  int
		ok   bool
	}{
		{"Floor hit", func() (int, bool) { return s.Floor(20) }, 20, true},
		{"Floor below", func() (int, bool) { return s.Floor(25) }, 20, true},
		{"Floor min", func() (int, bool) { return s.Floor(10) }, 10, true},
		{"Floor max", func() (int, bool) { return s.Floor(50) }, 50, true},
		{"Floor miss", func() (int, bool) { return s.Floor(5) }, 0, false},
		{"Ceiling hit", func() (int, bool) { return s.Ceiling(20) }, 20, true},
		{"Ceiling above", func() (int, bool) { return s.Ceiling(25) }, 30, true},
		{"Ceiling min", func() (int, bool) { return s.Ceiling(10) }, 10, true},
		{"Ceiling max", func() (int, bool) { return s.Ceiling(50) }, 50, true},
		{"Ceiling miss", func() (int, bool) { return s.Ceiling(60) }, 0, false},
		{"Lower hit", func() (int, bool) { return s.Lower(20) }, 10, true},
		{"Lower gap", func() (int, bool) { return s.Lower(25) }, 20, true},
		{"Lower miss", func() (int, bool) { return s.Lower(10) }, 0, false},
		{"Higher hit", func() (int, bool) { return s.Higher(20) }, 30, true},
		{"Higher gap", func() (int, bool) { return s.Higher(25) }, 30, true},
		{"Higher miss", func() (int, bool) { return s.Higher(50) }, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.fn()
			testutil.EqVal(t, ok, tc.ok)
			testutil.EqVal(t, got, tc.val)
		})
	}

	empty := New(testutil.LessInt)
	for _, tc := range []struct {
		name string
		fn   func() (int, bool)
	}{
		{"empty Floor", func() (int, bool) { return empty.Floor(1) }},
		{"empty Ceiling", func() (int, bool) { return empty.Ceiling(1) }},
		{"empty Lower", func() (int, bool) { return empty.Lower(1) }},
		{"empty Higher", func() (int, bool) { return empty.Higher(1) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.fn()
			testutil.EqVal(t, ok, false)
			testutil.EqVal(t, got, 0)
		})
	}
}

func TestRange(t *testing.T) {
	s := Of(testutil.LessInt, 10, 20, 30, 40, 50)

	cases := []struct {
		name   string
		lo, hi int
		want   []int
	}{
		{"inclusive", 20, 40, []int{20, 30, 40}},
		{"gap", 25, 35, []int{30}},
		{"single", 30, 30, []int{30}},
		{"full", 10, 50, []int{10, 20, 30, 40, 50}},
		{"miss", 100, 200, []int{}},
		{"lo>hi", 40, 20, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.Eq(t, collectRange(s, tc.lo, tc.hi), tc.want)
		})
	}

	t.Run("empty tree", func(t *testing.T) {
		testutil.Eq(t, collectRange(New(testutil.LessInt), 1, 9), []int{})
	})
	t.Run("early break", func(t *testing.T) {
		var got []int
		for v := range s.Range(20, 50) {
			got = append(got, v)
			if v == 30 {
				break
			}
		}
		testutil.Eq(t, got, []int{20, 30})
	})
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
			a, b := Of(testutil.LessInt, tc.a...), Of(testutil.LessInt, tc.b...)
			eqSet(t, a.Union(b), tc.want)
			eqSet(t, a, orEmpty(tc.a))
			eqSet(t, b, orEmpty(tc.b))
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
			a, b := Of(testutil.LessInt, tc.a...), Of(testutil.LessInt, tc.b...)
			eqSet(t, a.Intersect(b), tc.want)
			eqSet(t, a, orEmpty(tc.a))
			eqSet(t, b, orEmpty(tc.b))
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
			a, b := Of(testutil.LessInt, tc.a...), Of(testutil.LessInt, tc.b...)
			eqSet(t, a.Difference(b), tc.want)
			eqSet(t, a, orEmpty(tc.a))
			eqSet(t, b, orEmpty(tc.b))
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
			a, b := Of(testutil.LessInt, tc.a...), Of(testutil.LessInt, tc.b...)
			eqSet(t, a.SymmetricDifference(b), tc.want)
			eqSet(t, a, orEmpty(tc.a))
			eqSet(t, b, orEmpty(tc.b))
		})
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		name string
		a, b *OrderedSet[int]
		want bool
	}{
		{"same", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 2, 1), true},
		{"different", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 1, 2, 3), false},
		{"empty", New(testutil.LessInt), Of(testutil.LessInt), true},
		{"empty vs value", New(testutil.LessInt), Of(testutil.LessInt, 1), false},
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
		a, b *OrderedSet[int]
		want bool
	}{
		{"subset", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 1, 2, 3), true},
		{"not subset", Of(testutil.LessInt, 1, 2, 3), Of(testutil.LessInt, 1, 2), false},
		{"equal", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 2, 1), true},
		{"empty of values", New(testutil.LessInt), Of(testutil.LessInt, 1), true},
		{"empty of empty", New(testutil.LessInt), New(testutil.LessInt), true},
		{"value of empty", Of(testutil.LessInt, 1), New(testutil.LessInt), false},
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
		a, b *OrderedSet[int]
		want bool
	}{
		{"superset", Of(testutil.LessInt, 1, 2, 3), Of(testutil.LessInt, 1, 2), true},
		{"not superset", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 1, 2, 3), false},
		{"equal", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 2, 1), true},
		{"values of empty", Of(testutil.LessInt, 1), New(testutil.LessInt), true},
		{"empty of empty", New(testutil.LessInt), New(testutil.LessInt), true},
		{"empty of values", New(testutil.LessInt), Of(testutil.LessInt, 1), false},
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
		a, b *OrderedSet[int]
		want bool
	}{
		{"disjoint", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 3, 4), true},
		{"overlap", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 2, 3), false},
		{"equal", Of(testutil.LessInt, 1, 2), Of(testutil.LessInt, 1, 2), false},
		{"empty other", Of(testutil.LessInt, 1), New(testutil.LessInt), true},
		{"empty self", New(testutil.LessInt), Of(testutil.LessInt, 1), true},
		{"both empty", New(testutil.LessInt), New(testutil.LessInt), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.EqVal(t, tc.a.IsDisjoint(tc.b), tc.want)
		})
	}
}

func TestAlgebraUsesReceiverLess(t *testing.T) {
	a := Of(testutil.GreaterInt, 1, 3)
	b := Of(testutil.LessInt, 2, 4)
	got := a.Union(b)
	eqSet(t, got, []int{4, 3, 2, 1})
}

func TestNilOperand(t *testing.T) {
	eqSet(t, Of(testutil.LessInt, 1, 2).Union(nil), []int{1, 2})
	eqSet(t, Of(testutil.LessInt, 1).Intersect(nil), []int{})
	eqSet(t, Of(testutil.LessInt, 1, 2).Difference(nil), []int{1, 2})
	eqSet(t, Of(testutil.LessInt, 1, 2).SymmetricDifference(nil), []int{1, 2})
	testutil.EqVal(t, Of(testutil.LessInt, 1).Equal(nil), false)
	testutil.EqVal(t, New(testutil.LessInt).Equal(nil), true)
	testutil.EqVal(t, Of(testutil.LessInt, 1).IsSubsetOf(nil), false)
	testutil.EqVal(t, New(testutil.LessInt).IsSubsetOf(nil), true)
	testutil.EqVal(t, Of(testutil.LessInt, 1).IsSupersetOf(nil), true)
	testutil.EqVal(t, New(testutil.LessInt).IsSupersetOf(nil), true)
	testutil.EqVal(t, Of(testutil.LessInt, 1).IsDisjoint(nil), true)
}

// Nil less is deferred: empty New/Of succeed; panic on first compare.
func TestNilLess(t *testing.T) {
	s := New[int](nil)
	testutil.EqVal(t, s.Len(), 0)
	testutil.EqVal(t, Of[int](nil).Len(), 0)
	testutil.EqVal(t, Of[int](nil, 1).Len(), 1)

	single := New[int](nil)
	single.Add(1)
	testutil.EqVal(t, single.Len(), 1)
	first, ok := single.First()
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, first, 1)

	cases := []struct {
		name string
		fn   func()
	}{
		{"Of two", func() { Of[int](nil, 1, 2) }},
		{"FromSlice", func() { FromSlice[int](nil, []int{1, 2}) }},
		{"Add second", func() {
			s := New[int](nil)
			s.Add(1)
			s.Add(2)
		}},
		{"Contains after one", func() {
			s := New[int](nil)
			s.Add(1)
			s.Contains(1)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func TestNilReceiverPanics(t *testing.T) {
	var s *OrderedSet[int]
	cases := []struct {
		name string
		fn   func()
	}{
		{"Add", func() { s.Add(1) }},
		{"Remove", func() { s.Remove(1) }},
		{"Clear", func() { s.Clear() }},
		{"Len", func() { s.Len() }},
		{"Union", func() { s.Union(Of(testutil.LessInt, 1)) }},
		{"Equal", func() { s.Equal(New(testutil.LessInt)) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func TestAddRemoveShuffle(t *testing.T) {
	vals := []int{5, 1, 9, 3, 7, 2, 8, 4, 6, 0}
	s := New(testutil.LessInt)
	ref := map[int]struct{}{}

	for _, v := range vals {
		s.Add(v)
		ref[v] = struct{}{}
		assertMatchesRef(t, s, ref)
	}

	for _, v := range []int{3, 7, 0, 9, 1, 5, 2, 8, 4, 6} {
		testutil.EqVal(t, s.Remove(v), true)
		delete(ref, v)
		assertMatchesRef(t, s, ref)
	}
}

func TestAddRemoveRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	s := New(testutil.LessInt)
	ref := map[int]struct{}{}

	for range 200 {
		v := rng.Intn(50)
		if rng.Intn(2) == 0 || len(ref) == 0 {
			s.Add(v)
			ref[v] = struct{}{}
		} else {
			keys := sortedKeys(ref)
			v = keys[rng.Intn(len(keys))]
			testutil.EqVal(t, s.Remove(v), true)
			delete(ref, v)
		}
		assertMatchesRef(t, s, ref)
	}
}

func eqSet[T comparable](t *testing.T, s *OrderedSet[T], want []T) {
	t.Helper()
	got := s.ToSlice()
	testutil.Eq(t, got, want)
	if len(want) == 0 {
		testutil.NonNilEmpty(t, got)
	}
}

func collectRange[T comparable](s *OrderedSet[T], lo, hi T) []T {
	var out []T
	for v := range s.Range(lo, hi) {
		out = append(out, v)
	}
	if out == nil {
		return []T{}
	}
	return out
}

func orEmpty(vals []int) []int {
	if vals == nil {
		return []int{}
	}
	return vals
}

func assertMatchesRef(t *testing.T, s *OrderedSet[int], ref map[int]struct{}) {
	t.Helper()
	want := sortedKeys(ref)
	testutil.EqVal(t, s.Len(), len(want))
	testutil.Eq(t, s.ToSlice(), want)
	for _, v := range want {
		testutil.EqVal(t, s.Contains(v), true)
	}
}

func sortedKeys(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
