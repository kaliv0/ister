package rbtree

import (
	"iter"
)

type node[K any, V any] struct {
	key    K
	val    V
	parent *node[K, V]
	left   *node[K, V]
	right  *node[K, V]
	red    bool
}

type Tree[K any, V any] struct {
	less func(a, b K) bool
	root *node[K, V]
	size int
}

// Empty tree. Nil root until first Put.
// A nil less panics when ordering is needed (e.g. second Put), not in New.
// First Put on an empty tree does not call less.
func New[K any, V any](less func(a, b K) bool) *Tree[K, V] {
	return &Tree[K, V]{less: less}
}

// Number of keys.
func (t *Tree[K, V]) Len() int {
	return t.size
}

// Insert or replace. Returns previous value and whether the key was already present.
func (t *Tree[K, V]) Put(k K, v V) (V, bool) {
	n, parent := t.find(k)
	if n != nil {
		// found k -> set new v and return old one
		old := n.val
		n.val = v
		return old, true
	}
	var zero V
	// not found -> parent is attachment point
	if parent == nil {
		// empty tree -> first insert
		n := node[K, V]{key: k, val: v, red: false}
		t.root = &n
		t.size++
		return zero, false
	}
	// insert as leaf
	n = &node[K, V]{key: k, val: v, red: true, parent: parent}
	if t.less(k, parent.key) {
		parent.left = n
	} else {
		parent.right = n
	}
	t.size++

	// rb fixup
	t.insertFixup(n)
	return zero, false
}

func (t *Tree[K, V]) insertFixup(n *node[K, V]) {
	// while parent is red we have a red-red violation
	for n.parent != nil && n.parent.red {
		parent := n.parent
		grand := parent.parent // exists: root is always black
		if parent == grand.left {
			n = t.insertFixupLeft(n, parent, grand)
		} else {
			n = t.insertFixupRight(n, parent, grand)
		}
	}
	t.root.red = false // root is always black
}

// Parent is left child of grand. Uncle is on the right.
// Returns the node to continue fixup from.
func (t *Tree[K, V]) insertFixupLeft(n, parent, grand *node[K, V]) *node[K, V] {
	uncle := grand.right
	if isRed(uncle) {
		// Uncle red → recolor and move up.
		parent.red = false
		uncle.red = false
		grand.red = true
		return grand
	}
	// Uncle black (or nil) → rotate into outer case, then recolor.
	if n == parent.right {
		// Inner case: n is right child → make it outer.
		n = parent
		t.rotateLeft(n)
		parent = n.parent
	}
	parent.red = false
	grand.red = true
	t.rotateRight(grand)
	return n
}

// Mirror: parent is right child of grand. Uncle is on the left.
func (t *Tree[K, V]) insertFixupRight(n, parent, grand *node[K, V]) *node[K, V] {
	uncle := grand.left
	if isRed(uncle) {
		parent.red = false
		uncle.red = false
		grand.red = true
		return grand
	}
	if n == parent.left {
		n = parent
		t.rotateRight(n)
		parent = n.parent
	}
	parent.red = false
	grand.red = true
	t.rotateLeft(grand)
	return n
}

func (t *Tree[K, V]) rotateLeft(x *node[K, V]) {
	/*
		y takes x's place as subtree root

		      P                     P
		      |                     |
		      x                     y
		     / \                   / \
		   A    y        --->     x   C
		       / \               / \
		      B   C             A  B
	*/

	y := x.right

	// move y's left subtree under x as x's right child
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}

	// link y to x's original parent
	y.parent = x.parent
	if x.parent == nil {
		// x was root -> y becomes new root
		t.root = y
	} else if x == x.parent.left {
		// x was left child
		x.parent.left = y
	} else {
		// x was right child
		x.parent.right = y
	}

	// attach x under y as left child
	y.left = x
	x.parent = y
}

func (t *Tree[K, V]) rotateRight(x *node[K, V]) {
	/*
		y takes x's place as subtree root

		      P                     P
		      |                     |
		      x                     y
		     / \                   / \
		    y   C        --->     A   x
		   / \                       / \
		  A   B                     B   C
	*/

	y := x.left

	// move y's right subtree under x as x's left child
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}

	// link y to x's original parent
	y.parent = x.parent
	if x.parent == nil {
		// x was root -> y becomes new root
		t.root = y
	} else if x == x.parent.left {
		// x was left child
		x.parent.left = y
	} else {
		// x was right child
		x.parent.right = y
	}

	// attach x under y as right child
	y.right = x
	x.parent = y
}

// Value for k. Missing: (zero, false).
func (t *Tree[K, V]) Get(k K) (V, bool) {
	n, _ := t.find(k)
	if n == nil {
		var zero V
		return zero, false
	}
	return n.val, true
}

// Whether k is present.
func (t *Tree[K, V]) Contains(k K) bool {
	n, _ := t.find(k)
	return n != nil
}

// Remove k. Returns previous value and whether it was present.
func (t *Tree[K, V]) Delete(k K) (old V, ok bool) {
	if t.root == nil {
		return
	}
	n, _ := t.find(k)
	if n == nil {
		return
	}
	old = n.val
	t.deleteNode(n)
	t.size--
	return old, true
}

// Unlink n from the tree and restore red-black invariants.
//
// Two-child case: copy the successor's key/val into n, then unlink the
// successor instead (successor has no left child). That avoids moving
// n's left/right pointers onto another node.
func (t *Tree[K, V]) deleteNode(n *node[K, V]) {
	// If n has two children, delete its successor instead of n.
	// Successor is the minimum of n.right — always has nil left.
	if n.left != nil && n.right != nil {
		succ := t.minimum(n.right)
		n.key = succ.key
		n.val = succ.val
		n = succ
	}

	// Now n has at most one child. child may be nil if n is a leaf.
	var child *node[K, V]
	if n.left != nil {
		child = n.left
	} else {
		child = n.right
	}
	// Parent of the hole after transplant. Needed when child is nil (no sentinel).
	childParent := n.parent
	wasRed := n.red

	t.transplant(n, child)

	// Deleting a red node cannot break black-height. Black deletion may.
	if !wasRed {
		t.deleteFixup(child, childParent)
	}
}

// Leftmost node in the subtree rooted at n. n must be non-nil.
func (t *Tree[K, V]) minimum(n *node[K, V]) *node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// Rightmost node in the subtree rooted at n. n must be non-nil.
func (t *Tree[K, V]) maximum(n *node[K, V]) *node[K, V] {
	for n.right != nil {
		n = n.right
	}
	return n
}

// Replace subtree rooted at old with subtree rooted at repl.
// repl may be nil (old is spliced out with no replacement).
func (t *Tree[K, V]) transplant(old, repl *node[K, V]) {
	if old.parent == nil {
		t.root = repl
	} else if old == old.parent.left {
		old.parent.left = repl
	} else {
		old.parent.right = repl
	}
	if repl != nil {
		repl.parent = old.parent
	}
}

// Fix "double-black" at n after deleting a black node.
// n may be nil; then parent.left/right == nil tells which side the hole is on.
func (t *Tree[K, V]) deleteFixup(n, parent *node[K, V]) {
	for n != t.root && !isRed(n) {
		if parent.left == n {
			n, parent = t.deleteFixupLeft(n, parent)
		} else {
			n, parent = t.deleteFixupRight(n, parent)
		}
	}
	if n != nil {
		n.red = false
	}
}

// n is a left child (or nil in the left slot). Sibling is on the right.
// Returns the next (n, parent) for the fixup loop.
func (t *Tree[K, V]) deleteFixupLeft(n, parent *node[K, V]) (*node[K, V], *node[K, V]) {
	sibling := parent.right
	if isRed(sibling) {
		// Case 1: red sibling → recolor and rotate so sibling becomes black.
		sibling.red = false
		parent.red = true
		t.rotateLeft(parent)
		sibling = parent.right
	}
	if sibling == nil {
		// No sibling: move double-black up to parent.
		return parent, parent.parent
	}
	if !isRed(sibling.left) && !isRed(sibling.right) {
		// Case 2: sibling black, both nephews black → recolor sibling, move up.
		sibling.red = true
		return parent, parent.parent
	}
	if !isRed(sibling.right) {
		// Case 3: near nephew red, far nephew black → rotate sibling → case 4.
		if sibling.left != nil {
			sibling.left.red = false
		}
		sibling.red = true
		t.rotateRight(sibling)
		sibling = parent.right
	}
	// Case 4: far nephew red → recolor, rotate parent, done.
	sibling.red = parent.red
	parent.red = false
	if sibling.right != nil {
		sibling.right.red = false
	}
	t.rotateLeft(parent)
	return t.root, nil
}

// Mirror of deleteFixupLeft: n is a right child. Sibling is on the left.
func (t *Tree[K, V]) deleteFixupRight(n, parent *node[K, V]) (*node[K, V], *node[K, V]) {
	sibling := parent.left
	if isRed(sibling) {
		sibling.red = false
		parent.red = true
		t.rotateRight(parent)
		sibling = parent.left
	}
	if sibling == nil {
		// NB:  in a valid RB tree after a black delete this should not happen
		return parent, parent.parent
	}
	if !isRed(sibling.left) && !isRed(sibling.right) {
		sibling.red = true
		return parent, parent.parent
	}
	if !isRed(sibling.left) {
		if sibling.right != nil {
			sibling.right.red = false
		}
		sibling.red = true
		t.rotateLeft(sibling)
		sibling = parent.left
	}
	sibling.red = parent.red
	parent.red = false
	if sibling.left != nil {
		sibling.left.red = false
	}
	t.rotateRight(parent)
	return t.root, nil
}

func isRed[K any, V any](n *node[K, V]) bool {
	return n != nil && n.red
}

// Drop all entries. Keep less so later Puts can reuse the tree.
func (t *Tree[K, V]) Clear() {
	t.root = nil
	t.size = 0
}

//---- Navigable point queries ----//

// Smallest key. Empty: (zero, zero, false).
func (t *Tree[K, V]) Min() (key K, val V, ok bool) {
	if t.root == nil {
		return
	}
	n := t.minimum(t.root)
	return n.key, n.val, true
}

// Largest key. Empty: (zero, zero, false).
func (t *Tree[K, V]) Max() (key K, val V, ok bool) {
	if t.root == nil {
		return
	}
	n := t.maximum(t.root)
	return n.key, n.val, true
}

// Greatest key ≤ k. None: (zero, zero, false).
func (t *Tree[K, V]) Floor(k K) (K, V, bool) {
	return t.floorBound(k, false)
}

// Greatest key < k. None: (zero, zero, false).
func (t *Tree[K, V]) Lower(k K) (K, V, bool) {
	return t.floorBound(k, true)
}

// Least key ≥ k. None: (zero, zero, false).
func (t *Tree[K, V]) Ceiling(k K) (K, V, bool) {
	return t.ceilingBound(k, false)
}

// Least key > k. None: (zero, zero, false).
func (t *Tree[K, V]) Higher(k K) (K, V, bool) {
	return t.ceilingBound(k, true)
}

// strict false → Floor (≤ k), true → Lower (< k).
func (t *Tree[K, V]) floorBound(k K, strict bool) (key K, val V, ok bool) {
	if t.root == nil {
		return
	}

	var cand *node[K, V]
	for n := t.root; n != nil; {
		if t.less(k, n.key) {
			n = n.left
		} else if t.less(n.key, k) {
			cand = n
			n = n.right
		} else if strict {
			n = n.left
		} else {
			return n.key, n.val, true
		}
	}
	if cand != nil {
		return cand.key, cand.val, true
	}
	return
}

// strict false → Ceiling (≥ k), true → Higher (> k).
func (t *Tree[K, V]) ceilingBound(k K, strict bool) (key K, val V, ok bool) {
	if t.root == nil {
		return
	}

	var cand *node[K, V]
	for n := t.root; n != nil; {
		if t.less(k, n.key) {
			cand = n
			n = n.left
		} else if t.less(n.key, k) {
			n = n.right
		} else if strict {
			n = n.right
		} else {
			return n.key, n.val, true
		}
	}
	if cand != nil {
		return cand.key, cand.val, true
	}
	return
}

// Ascending key order. Early stop when yield returns false.
func (t *Tree[K, V]) Ascend() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var zero K
		t.walkAscend(t.root, false, zero, zero, yield)
	}
}

// Ascending keys in [lo, hi] inclusive. Early stop when yield returns false.
func (t *Tree[K, V]) AscendRange(lo, hi K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if t.less(hi, lo) {
			return
		}
		t.walkAscend(t.root, true, lo, hi, yield)
	}
}

func (t *Tree[K, V]) walkAscend(n *node[K, V], bounded bool, lo, hi K, yield func(K, V) bool) bool {
	if n == nil {
		return true
	}
	// prune/filter subtrees for range walk
	if bounded && t.less(n.key, lo) {
		return t.walkAscend(n.right, true, lo, hi, yield) // all left too small
	}
	if bounded && t.less(hi, n.key) {
		return t.walkAscend(n.left, true, lo, hi, yield) // all right too big
	}
	// in-order traversal left-node-right
	if !t.walkAscend(n.left, bounded, lo, hi, yield) {
		return false
	}
	if !yield(n.key, n.val) {
		return false
	}
	return t.walkAscend(n.right, bounded, lo, hi, yield)
}

func (t *Tree[K, V]) find(k K) (n, parent *node[K, V]) {
	for n = t.root; n != nil; {
		if t.less(k, n.key) {
			// k is smaller -> search in the left subtree
			parent = n
			n = n.left
		} else if t.less(n.key, k) {
			// k is bigger -> search in the right subtree
			parent = n
			n = n.right
		} else {
			// found k
			return n, parent
		}
	}
	return nil, parent
}
