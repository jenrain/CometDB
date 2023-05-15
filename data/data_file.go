package data

import "CometDB/fio"

// DataFile 数据文件
type DataFile struct {
	// 当前文件id
	FileId uint32

	// 文件写到了哪个位置
	WriteOff int64

	// io 读写管理
	IOManager fio.IOManager
}

// OpenDataFile 打开新的数据文件
func OpenDataFile(dirPath string, field uint32) (*DataFile, error) {
	return nil, nil
}

func (df *DataFile) ReadLogRecord(logOffset int64) (*LogRecord, error) {
	return nil, nil
}

func (df *DataFile) Write(buf []byte) error {
	return nil
}

// Sync 数据持久化到磁盘
func (df *DataFile) Sync() error {
	return nil
}
