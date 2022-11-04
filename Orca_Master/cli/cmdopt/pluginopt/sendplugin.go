package pluginopt

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/util"
	"encoding/json"
)

const SliceBytes = 40 * 1024 // 分片大小

type ShellcodeMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	Params     string // 参数
	SliceNum   int    // 基础分片数量
	SliceSize  int64  // 基础分片大小
	RemainSize int64  // 剩余分片大小
}

// 获取程序集元信息，并加密
func GetShellcodeMetaInfo(uploadFile, params string) string {
	sliceNum, remainSize := fileopt.GetFileSliceInfo(uploadFile)
	sliceSize := int64(SliceBytes)
	assemblyMetaInfo := ShellcodeMetaInfo{
		Fid:        util.GenUUID(),
		Params:     params,
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
func SendShellcodeMetaMsg(clientId, metaData string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "plugin"
	data := metaData
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}
