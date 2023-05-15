package CometDB

// Options 配置结构体
type Options struct {
	// 数据库数据目录
	DirPath string

	// 单个数据文件大小阈值
	DataFileSize int64

	// 每次写数据是否持久化
	SyncWrites bool
}
