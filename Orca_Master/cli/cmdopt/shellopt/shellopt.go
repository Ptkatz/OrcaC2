package shellopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
)

// 填充ClientInfo结构并发送execShell消息
func SendExecShellMsg(clientId, cmdStr, msgId string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "execShell"
	cmdInfo := &common.CmdInfo{
		Context: cmdStr,
		Attach:  "",
	}
	cmdData, _ := json.Marshal(cmdInfo)
	data, _ := crypto.Encrypt(cmdData, []byte(config.AesKey))
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data, msgId)
	return retData
}

// 打印命令回显
func PrintShellOutput(message string) {
	var retData common.RetData
	err := json.Unmarshal([]byte(message), &retData)
	if err != nil {
		return
	}
	data := retData.Data.(string)
	output, err := crypto.Decrypt(data, []byte(config.AesKey))
	if err != nil {
		return
	}
	fmt.Println(output)
}
