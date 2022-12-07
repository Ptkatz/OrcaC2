package rdesktopopt

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
	"github.com/go-vgo/robotgo"
)

type ScreenSize struct {
	X int
	Y int
}

// 发送截图数据
func SendScreenSize(clientId string) {
	x, y := robotgo.GetScaleSize()
	var screenSize = ScreenSize{
		X: x,
		Y: y,
	}
	screenSizeJson, _ := json.Marshal(screenSize)
	data, _ := crypto.Encrypt(screenSizeJson, []byte(config.AesKey))
	msg := "screenSize"
	sendUserId := common.ClientId

	common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
}
