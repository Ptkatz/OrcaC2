package controller

import (
	"Orca_Master/cli/cmdopt/listopt"
	"Orca_Master/cli/cmdopt/sshopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"github.com/desertbit/grumble"
	"time"
)

// 选择id对应的主机命令
var selectCmd = &grumble.Command{
	Name:  "select",
	Help:  "select the host id waiting to be operated",
	Usage: "select [-h | --help] <id>",
	Args: func(a *grumble.Args) {
		a.Int("id", "set host id")
	},
	Run: func(c *grumble.Context) error {
		SelectId = c.Args.Int("id")
		// 发送命令获取HostList
		common.SendSuccessMsg("Server", common.ClientId, "listHosts", "")
		// 从管道中接收消息
		select {
		case msg := <-common.DefaultMsgChan:
			HostLists = listopt.GetHostLists(msg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		// 判断设置id是否合法
		if SelectId > len(HostLists) {
			SelectId = -1
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "there is no corresponding host in the list")
			return nil
		}
		// 更改终端提示符
		host := HostLists[SelectId-1]
		SelectClientId = host.ClientId
		SelectIp = host.Ip
		SelectVer = host.Version
		sshopt.MySsh.Node = SelectClientId
		if SelectId <= 0 {
			sshopt.MySsh.Node = "Server"
		}
		App.SetPrompt("Orca[" + Uname + "] → " + SelectIp + " » ")
		RemoveCommand()
		AddCommand()
		sshopt.InitSshOption()
		return nil
	},
}
