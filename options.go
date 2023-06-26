package CometDB

import "os"

// Options 配置结构体
type Options struct {
	// 数据库数据目录
	DirPath string

	// 单个数据文件大小阈值
	DataFileSize int64

	// 每次写数据是否持久化
	SyncWrites bool

	// 索引类型
	IndexType IndexerType
}

// IteratorOptions 索引迭代器配置项
type IteratorOptions struct {

	// 遍历前缀为指定值的 key，默认为空
	Prefix []byte

	// 是否要反向遍历，默认false是正向
	Reverse bool
}

type IndexerType = int8

const (
	// Btree 索引
	Btree IndexerType = iota + 1

	// ART 自适应基数树索引
	ART
)

var DefaultOptions = Options{
	DirPath:      os.TempDir(),
	DataFileSize: 256 * 1024 * 1024,
	SyncWrites:   false,
	IndexType:    Btree,
}

// DefaultIteratorOptions 默认迭代器
var DefaultIteratorOptions = IteratorOptions{
	Prefix:  nil,
	Reverse: false,
}
