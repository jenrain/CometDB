package index

import (
	"CometDB/data"
	"bytes"
	"github.com/google/btree"
	"sort"
	"sync"
)

// BTree 索引，主要封装了google的btree库
type BTree struct {
	tree *btree.BTree
	lock *sync.RWMutex
}

// NewBTree 初始化BTree索引结构
func NewBTree() *BTree {
	return &BTree{
		// 叶子节点的数量
		tree: btree.New(32),
		lock: new(sync.RWMutex),
	}
}

func (bt *BTree) Iterator(reverse bool) Iterator {
	if bt.tree == nil {
		return nil
	}
	bt.lock.RLock()
	defer bt.lock.RUnlock()
	return newBTreeIterator(bt.tree, reverse)
}

func (bt *BTree) Put(key []byte, pos *data.LogRecordPos) bool {
	it := &Item{
		key: key,
		pos: pos,
	}
	bt.lock.Lock()
	defer bt.lock.Unlock()
	bt.tree.ReplaceOrInsert(it)
	return true
}

func (bt *BTree) Get(key []byte) *data.LogRecordPos {
	it := &Item{
		key: key,
	}
	btreeItem := bt.tree.Get(it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos
}

func (bt *BTree) Size() int {
	return bt.tree.Len()
}

func (bt *BTree) Delete(key []byte) bool {
	it := &Item{key: key}
	bt.lock.Lock()
	defer bt.lock.Unlock()
	oldItem := bt.tree.Delete(it)
	if oldItem == nil {
		return false
	}
	return true
}

// BTree 索引迭代器
type btreeIterator struct {
	// 当前遍历的下标位置
	currIndex int

	// 是否是反向遍历
	reverse bool

	// key + 位置索引信息
	values []*Item
}

func newBTreeIterator(tree *btree.BTree, reverse bool) *btreeIterator {
	var idx int
	values := make([]*Item, tree.Len())

	// 将所有的数据存放到数组中
	saveValues := func(it btree.Item) bool {
		values[idx] = it.(*Item)
		idx++
		return true
	}

	if reverse {
		tree.Descend(saveValues)
	} else {
		tree.Ascend(saveValues)
	}

	return &btreeIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}
}

func (b *btreeIterator) Rewind() {
	b.currIndex = 0
}

func (b *btreeIterator) Seek(key []byte) {
	if b.reverse {
		b.currIndex = sort.Search(len(b.values), func(i int) bool {
			return bytes.Compare(b.values[i].key, key) >= 0
		})
	} else {
		b.currIndex = sort.Search(len(b.values), func(i int) bool {
			return bytes.Compare(b.values[i].key, key) <= 0
		})
	}

}

func (b *btreeIterator) Next() {
	b.currIndex += 1
}

func (b *btreeIterator) Valid() bool {
	return b.currIndex < len(b.values)
}

func (b *btreeIterator) Key() []byte {
	return b.values[b.currIndex].key
}

func (b *btreeIterator) Value() *data.LogRecordPos {
	return b.values[b.currIndex].pos
}

func (b *btreeIterator) Close() {
	b.values = nil
}
