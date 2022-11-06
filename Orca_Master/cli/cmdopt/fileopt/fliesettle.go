package fileopt

import (
	"Orca_Master/define/colorcode"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

var FileDataChan chan string

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

// 获取文件md5值
func GetFileMd5Sum(fileName string) string {
	pFile, _ := os.Open(fileName)
	defer pFile.Close()
	md5a := md5.New()
	io.Copy(md5a, pFile)
	return hex.EncodeToString(md5a.Sum(nil))
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

// 判断本地文件和合法性
func IsLocalFileLegal(localFile string) bool {
	if localFile == "" {
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "please enter the file path")
		return false
	}
	//if !filepath.IsAbs(localFile) {
	//	colorcode.PrintMessage(colorcode.SIGN_ERROR, "wrong file path")
	//	return false
	//}
	if !IsFile(localFile) {
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "local file is not exist")
		return false
	}
	return true
}

// 判断远端文件和合法性
func IsRemoteFileLegal(remoteFile string) bool {
	if remoteFile == "" {
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "please enter the file path")
		return false
	}

	return true
}
