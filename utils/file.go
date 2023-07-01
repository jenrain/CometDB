package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// DirSize 获取一个目录的大小
func DirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// AvailableDiskSize 获取磁盘剩余可用空间大小
func AvailableDiskSize() (uint64, error) {
	wd, err := syscall.Getwd()
	if err != nil {
		return 0, err
	}
	var stat syscall.Statfs_t
	if err = syscall.Statfs(wd, &stat); err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}

// CopyDir 拷贝数据目录
func CopyDir(src, dest string, exclude []string) error {
	// 目标不存在就创建
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}
	}

	// 递归遍历源目录下的所有文件
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		// 取出路径中 src 之后的部分
		fileName := strings.Replace(path, src, "", 1)
		if fileName == "" {
			return nil
		}

		// 遍历排除的目录
		for _, e := range exclude {
			// 匹配是否包含排除的目录
			matched, err := filepath.Match(e, info.Name())
			if err != nil {
				return err
			}
			if matched {
				return nil
			}
		}

		// 如果是目录，就直接创建目录下的所有文件
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(dest, fileName), info.Mode())
		}

		data, err := os.ReadFile(filepath.Join(src, fileName))
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dest, fileName), data, info.Mode())
	})
}
