package fileopt

import (
	"os"
)

// FileMetadata 文件片元
type FileMetaInfo struct {
	Fid          string // 操作文件ID，随机生成的UUID
	SaveFileName string // 保存的文件名称
	SliceNum     int    // 基础分片数量
	SliceSize    int64  // 基础分片大小
	RemainSize   int64  // 剩余分片大小
	Md5sum       string // 文件md5值
}

// 请求文件
type RequestFile struct {
	SaveFileName string
	DestFileName string
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给文件是否存在
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}
