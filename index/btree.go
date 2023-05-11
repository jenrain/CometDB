package index

import (
	"CometDB/data"
	"github.com/google/btree"
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
