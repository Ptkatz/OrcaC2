package common

import (
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
)

// 获取ClientId
func GetClientId(message string) string {
	retData := RetData{}
	err := json.Unmarshal([]byte(message), &retData)
	if err != nil {
		return ""
	}
	clientId := retData.Data.(map[string]interface{})["clientId"].(string)
	return clientId
}

func GetMsgId(message string) string {
	retData := RetData{}
	err := json.Unmarshal([]byte(message), &retData)
	if err != nil {
		return ""
	}
	return retData.MessageId
}

func SettleRetData(message string) (string, string, string) {
	msg, sendUserId, data := SettleRetDataNotDec(message)
	decData, _ := crypto.Decrypt(data, []byte(config.AesKey))
	return msg, sendUserId, decData
}

func SettleRetDataBt(message string) (string, string, []byte) {
	msg, sendUserId, data := SettleRetDataNotDec(message)
	decData, _ := crypto.DecryptBt(data, []byte(config.AesKey))
	return msg, sendUserId, decData
}

func SettleRetDataNotDec(message string) (string, string, string) {
	var retData = RetData{}
	err := json.Unmarshal([]byte(message), &retData)
	if err != nil {
		return "", "", ""
	}
	msg := retData.Msg
	decMsg, _ := crypto.Decrypt(msg, []byte(config.AesKey))
	sendUserId := retData.SendUserId
	data := retData.Data.(string)
	return decMsg, sendUserId, data
}
