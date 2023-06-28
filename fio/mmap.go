package fio

import (
	"golang.org/x/exp/mmap"
	"os"
)

// MMap IO，内存文件映射
type MMap struct {
	readerAt *mmap.ReaderAt
}

// NewMMapIOManager 初始化内存文件映射 IO
func NewMMapIOManager(fileName string) (*MMap, error) {
	_, err := os.OpenFile(fileName, os.O_CREATE, DataFilePerm)
	if err != nil {
		return nil, err
	}
	readerAt, err := mmap.Open(fileName)
	if err != nil {
		return nil, err
	}
	return &MMap{readerAt: readerAt}, nil

}

// Read 从文件的给定位置读取对应的数据
func (m *MMap) Read(b []byte, offset int64) (int, error) {
	return m.readerAt.ReadAt(b, offset)
}

// Write 写入字节数组到文件中
func (m *MMap) Write([]byte) (int, error) {
	panic("not implemented")
}

// Sync 持久化数据
func (m *MMap) Sync() error {
	panic("not implemented")
}

// Close 关闭文件
func (m *MMap) Close() error {
	return m.readerAt.Close()
}

// Size 获取到文件大小
func (m *MMap) Size() (int64, error) {
	return int64(m.readerAt.Len()), nil
}
