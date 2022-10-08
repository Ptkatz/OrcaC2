package smbopt

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/util"
	"encoding/json"
)

type SmbUploadStruct struct {
	SmbStruct    SmbOption
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
	msUpload := SmbUploadStruct{
		SmbStruct:    MySmb,
		FileMetaInfo: fileMetaInfo,
	}
	metaInfo, err := json.Marshal(msUpload)
	if err != nil {
		return ""
	}
	data, _ := crypto.Encrypt(metaInfo, []byte(config.AesKey))
	return data
}
