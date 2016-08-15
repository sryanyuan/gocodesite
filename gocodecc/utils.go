package gocodecc

import (
	"os"
)

// 获取文件大小的接口
type FileSize interface {
	Size() int64
}

// 获取文件信息的接口
type FileStat interface {
	Stat() (os.FileInfo, error)
}
