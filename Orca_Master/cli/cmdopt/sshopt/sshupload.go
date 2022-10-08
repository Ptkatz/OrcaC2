package sshopt

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/util"
	"encoding/json"
)

type SshUploadStruct struct {
	SshStruct    SshOption
	FileMetaInfo fileopt.FileMetaInfo
}

const SliceBytes = 4 * 1024 // 分片大小

// 获取文件元信息，并加密
func GetFileMetaInfo(uploadFile, saveFile string) string {
	sliceNum, remainSize := fileopt.GetFileSliceInfo(uploadFile)
	sliceSize := int64(SliceBytes)
	md5sum := fileopt.GetFileMd5Sum(uploadFile)
	fileMetaInfo := fileopt.FileMetaInfo{
		Fid:          util.GenUUID(),
		SaveFileName: saveFile,
		SliceNum:     sliceNum,
		SliceSize:    sliceSize,
		RemainSize:   remainSize,
		Md5sum:       md5sum,
	}
	sshUpload := SshUploadStruct{
		SshStruct:    MySsh,
		FileMetaInfo: fileMetaInfo,
	}
	metaInfo, err := json.Marshal(sshUpload)
	if err != nil {
		return ""
	}
	data, _ := crypto.Encrypt(metaInfo, []byte(config.AesKey))
	return data
}
