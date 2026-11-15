package backend

import (
	"github.com/google/btree"
)

// IndexItem pairs the integer index (ID) with the generic table value (T).
// The less function used by the B-Tree will sort items based exclusively on ID.
type IndexItem[T any] struct {
	// Now there is no constraint to validate ID uniqueness. Is is assumed for simplicity.
	ID    int
	value T
}

type BTreeIndex[T any] struct {
	tree *btree.BTreeG[IndexItem[T]]
}

// NewBtreeIndex initializes the B-Tree.
//
// - How to choose/calculate BTree degree:
//   - Lower degrees mean smaller nodes, leading to shallower memory overhead for small datasets or frequent pointer mutations.
//   - Higher degrees mean wider nodes with fewer pointer hops, which reduces cache misses on modern CPU architectures.
//
// - Performance Impact:
//   - Too low: The tree becomes too tall, requiring more node-traversal steps (though each step is fast).
//   - Too high: Traversing *inside* a single node requires a linear or binary search over a large array,
//     which can slow down insertion and search times if the node gets too wide.
func NewBtreeIndex[T any]() *BTreeIndex[T] {
	tree := btree.NewG(3, func(a, b IndexItem[T]) bool {
		return a.ID < b.ID
	})
	return &BTreeIndex[T]{tree: tree}
}

func (t *BTreeIndex[T]) Search(id int) (_ T, _ bool) {
	if item, found := t.tree.Get(IndexItem[T]{ID: id}); found {
		return item.value, true
	}

	return
}

// TODO: There is no concurrency guard for simplicity. In a production env it needs a mutex to
// make writes thread-safe
func (t *BTreeIndex[T]) Insert(id int, value T) (_ T, _ bool) {
	indexItem := IndexItem[T]{ID: id, value: value}
	if item, added := t.tree.ReplaceOrInsert(indexItem); added {
		return item.value, true
	}
	return
}
