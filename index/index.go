package index

import (
	"CometDB/data"
	"bytes"
	"github.com/google/btree"
)

// Indexer 抽象索引接口，后续如果想要接入其他的数据结构，直接实现这个接口即可
type Indexer interface {
	// Put 向索引中存储key对应答的数据位置信息
	Put(key []byte, pos *data.LogRecordPos) bool

	// Get 根据key取出对应的索引位置信息
	Get(key []byte) *data.LogRecordPos

	// Delete 根据key删除对应的索引位置信息
	Delete(key []byte) bool
}

type IndexType = int8

const (
	// Btree 索引
	Btree IndexType = iota + 1

	// ART 自适应基数树索引
	ART
)

// NewIndexer 根据类型初始化索引
func NewIndexer(typ IndexType) Indexer {
	switch typ {
	case Btree:
		return NewBTree()
	case ART:
		return nil
	default:
		panic("unsupported index type")
	}
}

// Item 实现btree中的Item接口，b树中实际存储的对象
type Item struct {
	key []byte
	pos *data.LogRecordPos
}

// Less 比较方法
func (ai Item) Less(bi btree.Item) bool {
	// a < b返回-1
	return bytes.Compare(ai.key, bi.(*Item).key) == -1
}
