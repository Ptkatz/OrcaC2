package common

import (
	"Orca_Puppet/define/api"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/retcode"
	"Orca_Puppet/tools/crypto"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// 发送消息
func SendToClient(url string, payload []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("systemId", config.SystemId)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func SendMsg(clientId, sendUserId, msg, data, msgId string, code int) HttpRetData {
	encMsg, _ := crypto.Encrypt([]byte(msg), []byte(config.AesKey))
	clientInfo := &ClientInfo{
		ClientId:   clientId,
		SendUserId: sendUserId,
		MessageId:  msgId,
		Code:       code,
		Msg:        encMsg,
		Data:       &data,
	}
	message, _ := json.Marshal(clientInfo)
	body := SendToClient(api.SEND_TO_CLIENT_API, message)
	var retData HttpRetData
	err := json.Unmarshal(body, &retData)
	if err != nil {
		return HttpRetData{retcode.FAIL, "", ""}
	}
	return retData
}

func SendSuccessMsg(clientId, sendUserId, msg, data, msgId string) HttpRetData {
	return SendMsg(clientId, sendUserId, msg, data, msgId, retcode.SUCCESS)
}

func SendFailMsg(clientId, sendUserId, msg, data, msgId string) HttpRetData {
	return SendMsg(clientId, sendUserId, msg, data, msgId, retcode.FAIL)
}
