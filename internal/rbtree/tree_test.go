package rbtree

import (
	"maps"
	"math/rand"
	"slices"
	"testing"

	"github.com/kaliv0/ister/internal/testutil"
)

func lessInt(a, b int) bool { return a < b }

func TestNew(t *testing.T) {
	tr := New[int, string](lessInt)
	testutil.EqVal(t, tr.Len(), 0)
	checkRB(t, tr)
}

func TestPutGetContains(t *testing.T) {
	tr := New[int, string](lessInt)

	old, had := tr.Put(10, "ten")
	testutil.EqVal(t, had, false)
	testutil.EqVal(t, old, "")
	testutil.EqVal(t, tr.Len(), 1)
	mustGet(t, tr, 10, "ten")
	testutil.EqVal(t, tr.Contains(10), true)
	testutil.EqVal(t, tr.Contains(99), false)
	testutil.EqVal(t, New[int, string](lessInt).Contains(1), false)

	old, had = tr.Put(10, "TEN")
	testutil.EqVal(t, had, true)
	testutil.EqVal(t, old, "ten")
	testutil.EqVal(t, tr.Len(), 1)
	mustGet(t, tr, 10, "TEN")
}

func TestGet(t *testing.T) {
	tr := New[int, string](lessInt)
	for _, e := range []struct {
		k int
		v string
	}{
		{20, "twenty"}, {10, "ten"}, {30, "thirty"}, {5, "five"}, {15, "fifteen"},
	} {
		tr.Put(e.k, e.v)
	}

	cases := []struct {
		name string
		key  int
		val  string
		ok   bool
	}{
		{"root", 20, "twenty", true},
		{"left", 10, "ten", true},
		{"leaf", 15, "fifteen", true},
		{"miss low", 1, "", false},
		{"miss high", 99, "", false},
		{"miss gap", 12, "", false},
		{"empty key", 0, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, ok := tr.Get(tc.key)
			testutil.EqVal(t, ok, tc.ok)
			testutil.EqVal(t, v, tc.val)
		})
	}

	t.Run("empty", func(t *testing.T) {
		v, ok := New[int, string](lessInt).Get(1)
		testutil.EqVal(t, ok, false)
		testutil.EqVal(t, v, "")
	})
}

func TestPutOrder(t *testing.T) {
	cases := []struct {
		name string
		keys []int
	}{
		{"ascending", []int{1, 2, 3, 4, 5, 6, 7}},
		{"descending", []int{7, 6, 5, 4, 3, 2, 1}},
		{"mixed", []int{30, 10, 20, 50, 40}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := New[int, int](lessInt)
			for _, k := range tc.keys {
				tr.Put(k, k*10)
			}
			want := slices.Clone(tc.keys)
			slices.Sort(want)
			testutil.EqVal(t, tr.Len(), len(want))
			testutil.Eq(t, keysOf(tr), want)
			for _, k := range want {
				mustGet(t, tr, k, k*10)
			}
			checkRB(t, tr)
		})
	}
}

func TestDelete(t *testing.T) {
	tr := New[int, string](lessInt)
	tr.Put(10, "ten")
	tr.Put(20, "twenty")

	old, ok := tr.Delete(10)
	testutil.EqVal(t, ok, true)
	testutil.EqVal(t, old, "ten")
	testutil.EqVal(t, tr.Len(), 1)
	testutil.EqVal(t, tr.Contains(10), false)
	checkRB(t, tr)

	old, ok = tr.Delete(10)
	testutil.EqVal(t, ok, false)
	testutil.EqVal(t, old, "")
}

func TestDeleteShapes(t *testing.T) {
	cases := []struct {
		name    string
		insert  []int
		deletes []int
	}{
		{"sole root", []int{1}, []int{1}},
		{"leaf", []int{2, 1, 3}, []int{1}},
		{"one child", []int{2, 1, 3, 4}, []int{3}},
		{"two children", []int{4, 2, 6, 1, 3, 5, 7}, []int{4}},
		{"until empty", []int{5, 3, 7, 1, 4, 6, 8}, []int{5, 3, 7, 1, 4, 6, 8}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := New[int, int](lessInt)
			ref := make(map[int]int)
			for _, k := range tc.insert {
				tr.Put(k, k*10)
				ref[k] = k * 10
			}
			checkRB(t, tr)
			for _, k := range tc.deletes {
				got, ok := tr.Delete(k)
				testutil.EqVal(t, ok, true)
				testutil.EqVal(t, got, ref[k])
				delete(ref, k)
				assertMatchesRef(t, tr, ref)
				checkRB(t, tr)
			}
			testutil.EqVal(t, tr.Len(), len(ref))
		})
	}
}

func TestAscend(t *testing.T) {
	tr := New[int, int](lessInt)
	for _, k := range []int{30, 10, 20, 50, 40} {
		tr.Put(k, k*10)
	}
	testutil.Eq(t, keysOf(tr), []int{10, 20, 30, 40, 50})

	var pairs [][2]int
	for k, v := range tr.Ascend() {
		pairs = append(pairs, [2]int{k, v})
	}
	testutil.EqVal(t, len(pairs), 5)
	for _, p := range pairs {
		testutil.EqVal(t, p[1], p[0]*10)
	}

	t.Run("early break", func(t *testing.T) {
		var got []int
		for k := range tr.Ascend() {
			got = append(got, k)
			if k == 30 {
				break
			}
		}
		testutil.Eq(t, got, []int{10, 20, 30})
	})
	t.Run("empty", func(t *testing.T) {
		for range New[int, int](lessInt).Ascend() {
			t.Fatal("empty Ascend should not yield")
		}
	})
}

func TestAscendRange(t *testing.T) {
	tr := New[int, int](lessInt)
	for _, k := range []int{10, 20, 30, 40, 50} {
		tr.Put(k, k)
	}

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
			testutil.Eq(t, keysOfRange(tr, tc.lo, tc.hi), tc.want)
		})
	}

	t.Run("empty tree", func(t *testing.T) {
		testutil.Eq(t, keysOfRange(New[int, int](lessInt), 1, 9), []int{})
	})
	t.Run("early break", func(t *testing.T) {
		var got []int
		for k := range tr.AscendRange(20, 50) {
			got = append(got, k)
			if k == 30 {
				break
			}
		}
		testutil.Eq(t, got, []int{20, 30})
	})
}

func TestNavigable(t *testing.T) {
	tr := New[int, int](lessInt)
	for _, k := range []int{10, 20, 30, 40, 50} {
		tr.Put(k, k*10)
	}

	cases := []struct {
		name string
		fn   func() (int, int, bool)
		key  int
		val  int
		ok   bool
	}{
		{"Min", func() (int, int, bool) { return tr.Min() }, 10, 100, true},
		{"Max", func() (int, int, bool) { return tr.Max() }, 50, 500, true},
		{"Floor hit", func() (int, int, bool) { return tr.Floor(20) }, 20, 200, true},
		{"Floor below", func() (int, int, bool) { return tr.Floor(25) }, 20, 200, true},
		{"Floor min", func() (int, int, bool) { return tr.Floor(10) }, 10, 100, true},
		{"Floor max", func() (int, int, bool) { return tr.Floor(50) }, 50, 500, true},
		{"Floor miss", func() (int, int, bool) { return tr.Floor(5) }, 0, 0, false},
		{"Ceiling hit", func() (int, int, bool) { return tr.Ceiling(20) }, 20, 200, true},
		{"Ceiling above", func() (int, int, bool) { return tr.Ceiling(25) }, 30, 300, true},
		{"Ceiling min", func() (int, int, bool) { return tr.Ceiling(10) }, 10, 100, true},
		{"Ceiling max", func() (int, int, bool) { return tr.Ceiling(50) }, 50, 500, true},
		{"Ceiling miss", func() (int, int, bool) { return tr.Ceiling(60) }, 0, 0, false},
		{"Lower hit", func() (int, int, bool) { return tr.Lower(20) }, 10, 100, true},
		{"Lower gap", func() (int, int, bool) { return tr.Lower(25) }, 20, 200, true},
		{"Lower miss", func() (int, int, bool) { return tr.Lower(10) }, 0, 0, false},
		{"Higher hit", func() (int, int, bool) { return tr.Higher(20) }, 30, 300, true},
		{"Higher gap", func() (int, int, bool) { return tr.Higher(25) }, 30, 300, true},
		{"Higher miss", func() (int, int, bool) { return tr.Higher(50) }, 0, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k, v, ok := tc.fn()
			testutil.EqVal(t, k, tc.key)
			testutil.EqVal(t, v, tc.val)
			testutil.EqVal(t, ok, tc.ok)
		})
	}

	empty := New[int, int](lessInt)
	for _, tc := range []struct {
		name string
		fn   func() (int, int, bool)
	}{
		{"empty Min", empty.Min},
		{"empty Max", empty.Max},
		{"empty Floor", func() (int, int, bool) { return empty.Floor(1) }},
		{"empty Ceiling", func() (int, int, bool) { return empty.Ceiling(1) }},
		{"empty Lower", func() (int, int, bool) { return empty.Lower(1) }},
		{"empty Higher", func() (int, int, bool) { return empty.Higher(1) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k, v, ok := tc.fn()
			testutil.EqVal(t, ok, false)
			testutil.EqVal(t, k, 0)
			testutil.EqVal(t, v, 0)
		})
	}
}

func TestClear(t *testing.T) {
	tr := New[int, string](lessInt)
	tr.Put(1, "a")
	tr.Put(2, "b")
	tr.Clear()
	testutil.EqVal(t, tr.Len(), 0)
	testutil.EqVal(t, tr.Contains(1), false)
	testutil.Eq(t, keysOf(tr), []int{})
	checkRB(t, tr)

	// less must still work: second Put compares.
	tr.Put(3, "c")
	tr.Put(1, "a")
	mustGet(t, tr, 3, "c")
	mustGet(t, tr, 1, "a")
	testutil.Eq(t, keysOf(tr), []int{1, 3})
	checkRB(t, tr)
}

func TestLessEquality(t *testing.T) {
	tr := New[testutil.Pair, string](
		func(a, b testutil.Pair) bool { return a.A < b.A },
	)
	tr.Put(testutil.Pair{A: 1, B: 10}, "first")
	tr.Put(testutil.Pair{A: 1, B: 20}, "second")
	testutil.EqVal(t, tr.Len(), 1)
	mustGet(t, tr, testutil.Pair{A: 1, B: 10}, "second")
	mustGet(t, tr, testutil.Pair{A: 1, B: 99}, "second")
	testutil.EqVal(t, tr.Contains(testutil.Pair{A: 2, B: 1}), false)
}

// Nil less is deferred: empty New succeeds; panic on first compare (needs a second key).
func TestNilLess(t *testing.T) {
	tr := New[int, string](nil)
	testutil.EqVal(t, tr.Len(), 0)

	_, ok := tr.Get(1)
	testutil.EqVal(t, ok, false)
	testutil.EqVal(t, tr.Contains(1), false)

	_, deleted := tr.Delete(1)
	testutil.EqVal(t, deleted, false)

	_, _, hasMin := tr.Min()
	testutil.EqVal(t, hasMin, false)

	for range tr.Ascend() {
		t.Fatal("empty Ascend should not yield")
	}

	// Single Put: no compare, same as PQ Push of one element.
	tr.Put(1, "a")
	testutil.EqVal(t, tr.Len(), 1)

	// Any further ordered op calls less (find on Put/Get, etc.).
	testutil.MustPanic(t, func() { tr.Put(2, "b") })
}

func TestNilReceiverPanics(t *testing.T) {
	var tr *Tree[int, string]
	cases := []struct {
		name string
		fn   func()
	}{
		{"Put", func() { tr.Put(1, "a") }},
		{"Get", func() { tr.Get(1) }},
		{"Contains", func() { tr.Contains(1) }},
		{"Delete", func() { tr.Delete(1) }},
		{"Clear", func() { tr.Clear() }},
		{"Len", func() { tr.Len() }},
		{"Min", func() { tr.Min() }},
		{"Max", func() { tr.Max() }},
		{"Floor", func() { tr.Floor(1) }},
		{"Ceiling", func() { tr.Ceiling(1) }},
		{"Lower", func() { tr.Lower(1) }},
		{"Higher", func() { tr.Higher(1) }},
		{"Ascend", func() {
			for range tr.Ascend() {
			}
		}},
		{"AscendRange", func() {
			for range tr.AscendRange(1, 2) {
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.MustPanic(t, tc.fn)
		})
	}
}

func TestPutDeleteShuffle(t *testing.T) {
	keys := []int{5, 1, 9, 3, 7, 2, 8, 4, 6, 0}
	tr := New[int, string](lessInt)
	ref := make(map[int]string)

	for _, k := range keys {
		ref[k] = string(rune('a' + k))
		tr.Put(k, ref[k])
		assertMatchesRef(t, tr, ref)
		checkRB(t, tr)
	}

	for _, k := range []int{3, 7, 0, 9, 1, 5, 2, 8, 4, 6} {
		got, ok := tr.Delete(k)
		testutil.EqVal(t, ok, true)
		testutil.EqVal(t, got, ref[k])
		delete(ref, k)
		assertMatchesRef(t, tr, ref)
		checkRB(t, tr)
	}
}

func TestPutDeleteRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	tr := New[int, int](lessInt)
	ref := make(map[int]int)

	for range 200 {
		k := rng.Intn(50)
		if rng.Intn(2) == 0 || len(ref) == 0 {
			tr.Put(k, k*10)
			ref[k] = k * 10
		} else {
			keys := sortedKeysAny(ref)
			k = keys[rng.Intn(len(keys))]
			got, ok := tr.Delete(k)
			testutil.EqVal(t, ok, true)
			testutil.EqVal(t, got, ref[k])
			delete(ref, k)
		}
		assertMatchesRef(t, tr, ref)
		checkRB(t, tr)
	}
}

func mustGet[K, V comparable](t *testing.T, tr *Tree[K, V], k K, want V) {
	t.Helper()
	got, ok := tr.Get(k)
	if !ok {
		t.Fatalf("Get(%v): missing", k)
	}
	testutil.EqVal(t, got, want)
}

func assertMatchesRef[V comparable](t *testing.T, tr *Tree[int, V], ref map[int]V) {
	t.Helper()
	testutil.EqVal(t, tr.Len(), len(ref))
	testutil.Eq(t, keysOf(tr), sortedKeysAny(ref))
	for k, want := range ref {
		mustGet(t, tr, k, want)
	}
}

func keysOf[K, V any](tr *Tree[K, V]) []K {
	var keys []K
	for k := range tr.Ascend() {
		keys = append(keys, k)
	}
	if keys == nil {
		return []K{}
	}
	return keys
}

func keysOfRange[K, V any](tr *Tree[K, V], lo, hi K) []K {
	var keys []K
	for k := range tr.AscendRange(lo, hi) {
		keys = append(keys, k)
	}
	if keys == nil {
		return []K{}
	}
	return keys
}

func sortedKeysAny[V any](m map[int]V) []int {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return keys
}

// checkRB verifies red-black invariants and that size matches the node count.
func checkRB[K, V any](t *testing.T, tr *Tree[K, V]) {
	t.Helper()
	if tr.root == nil {
		testutil.EqVal(t, tr.size, 0)
		return
	}
	if tr.root.red {
		t.Fatal("root must be black")
	}
	if tr.root.parent != nil {
		t.Fatal("root parent must be nil")
	}

	count := 0
	var walk func(n *node[K, V], blacks int) int
	walk = func(n *node[K, V], blacks int) int {
		if n == nil {
			return blacks
		}
		count++
		if n.red {
			if isRed(n.left) || isRed(n.right) {
				t.Fatalf("red-red at key %v", n.key)
			}
		} else {
			blacks++
		}
		if n.left != nil && n.left.parent != n {
			t.Fatalf("bad left parent at key %v", n.key)
		}
		if n.right != nil && n.right.parent != n {
			t.Fatalf("bad right parent at key %v", n.key)
		}
		if n.left != nil && !tr.less(n.left.key, n.key) {
			t.Fatalf("BST broken left of %v", n.key)
		}
		if n.right != nil && !tr.less(n.key, n.right.key) {
			t.Fatalf("BST broken right of %v", n.key)
		}
		lb := walk(n.left, blacks)
		rb := walk(n.right, blacks)
		if lb != rb {
			t.Fatalf("black-height mismatch at key %v: left=%d right=%d", n.key, lb, rb)
		}
		return lb
	}
	walk(tr.root, 0)
	testutil.EqVal(t, tr.size, count)
}
