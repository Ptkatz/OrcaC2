package controller

import (
	"Orca_Master/cli/cmdopt/listopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"github.com/desertbit/grumble"
	"time"
)

// 列出主机列表命令
var listCmd = &grumble.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Help:    "list hosts",
	Usage:   "list [-h | --help] [-i | --id id]",
	Flags: func(f *grumble.Flags) {
		f.Int("i", "id", -1, "view selected host information by id")
	},

	Run: func(c *grumble.Context) error {
		identity := c.Flags.Int("id")

		// 发送命令获取HostLists消息
		common.SendSuccessMsg("Server", common.ClientId, "listHosts", "")
		// 从管道中接收消息
		select {
		case msg := <-common.DefaultMsgChan:
			HostLists = listopt.GetHostLists(msg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		// 处理消息打印表格
		listopt.PrintTable(HostLists, identity)
		return nil
	},
}
