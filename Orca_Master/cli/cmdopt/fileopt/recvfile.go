package fileopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/json"
)

// 发送文件下载消息
func SendDownloadRequestMsg(clientId, sendUserId, destFile, saveFile string) common.HttpRetData {
	msg := "fileDownload"
	requestFile := RequestFile{DestFileName: destFile, SaveFileName: saveFile}
	requestFileData, _ := json.Marshal(requestFile)
	data, _ := crypto.Encrypt(requestFileData, []byte(config.AesKey))
	return common.SendSuccessMsg(clientId, sendUserId, msg, data)
}

// 从文件元消息中获取文件元
func GetMetaInfo(metaInfoMsg string) FileMetaInfo {
	_, _, decData := common.SettleRetData(metaInfoMsg)

	var fileMetaInfo FileMetaInfo
	json.Unmarshal([]byte(decData), &fileMetaInfo)
	return fileMetaInfo
}
