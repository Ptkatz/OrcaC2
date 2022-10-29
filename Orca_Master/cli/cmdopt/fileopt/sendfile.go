package fileopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/util"
	"encoding/json"
	"fmt"
	"github.com/tj/go-spin"
	"io"
	"math"
	"os"
	"time"
)

const SliceBytes = 40 * 1024 // 分片大小

// 获取文件大小
func GetFileSize(path string) int64 {
	fh, err := os.Stat(path)
	if err != nil {
		message := fmt.Sprintf("read file %s fail, err: %s\n", path, err)
		colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
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
	md5sum := GetFileMd5Sum(uploadFile)
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
func SendFileMetaMsg(clientId, metaData, msg string) common.HttpRetData {
	sendUserId := common.ClientId
	data := metaData
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

// 发送文件分片数据
func SendFileSliceMsg(clientId string, sliceData []byte) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "sliceData"
	data, _ := crypto.Encrypt(sliceData, []byte(config.AesKey))
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

// 发送文件
func SendFileData(clientId string, uploadFile string) {
	sliceNum, remainSize := GetFileSliceInfo(uploadFile)
	sliceSize := int64(SliceBytes)
	pUploadFile, _ := os.Open(uploadFile)
	defer pUploadFile.Close()
	s := spin.New()
	s.Set(spin.Box1)
	currentTime := time.Now().Format("2006/01/02 15:04:05")
	timeSign := colorcode.COLOR_GREY + currentTime + colorcode.END
	sign := fmt.Sprintf("%s %s", timeSign, colorcode.SIGN_NOTICE)
	for i := 0; i < sliceNum; i++ {
		fmt.Printf("\r%s%s uploading\033[m %s ", sign, colorcode.COLOR_CYAN, s.Next())
		sliceData := make([]byte, sliceSize)
		_, err := pUploadFile.Read(sliceData)
		if err != nil && err != io.EOF {
			panic(err.Error())
		}
		SendFileSliceMsg(clientId, sliceData)
	}
	// 处理最后一个分片
	sliceData := make([]byte, remainSize)
	_, err := pUploadFile.Read(sliceData)
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
	SendFileSliceMsg(clientId, sliceData)
	fmt.Println()
}
