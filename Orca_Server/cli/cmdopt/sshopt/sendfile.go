package sshopt

import (
	"Orca_Server/tools/util"
	"math"
	"os"
)

const SliceBytes = 4 * 1024 // 分片大小

// 获取文件片元信息
func GetFileSliceInfo(fileInfo os.FileInfo) (int, int64) {
	filesize := fileInfo.Size()                                           // 文件大小
	sliceBytes := SliceBytes                                              // 分片大小
	sliceNum := int(math.Ceil(float64(filesize)/float64(sliceBytes))) - 1 // 分片数量-1
	remainSize := filesize - int64(sliceNum*sliceBytes)                   // 最后一个分片大小
	return sliceNum, remainSize
}

// 获取文件元信息，并加密
func GetFileMetaInfo(remoteFileInfo os.FileInfo, saveFile string) FileMetaInfo {
	sliceNum, remainSize := GetFileSliceInfo(remoteFileInfo)
	sliceSize := int64(SliceBytes)
	fileMetaInfo := FileMetaInfo{
		Fid:          util.GenUUID(),
		SaveFileName: saveFile,
		SliceNum:     sliceNum,
		SliceSize:    sliceSize,
		RemainSize:   remainSize,
	}
	return fileMetaInfo
}
