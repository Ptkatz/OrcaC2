package controller

import (
	"Orca_Master/cli/cmdopt/shellopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/util"
	"github.com/desertbit/grumble"
	"path/filepath"
	"strings"
	"time"
)

// 发送Shell命令
var shellCmd = &grumble.Command{
	Name:    "shell",
	Aliases: []string{"sh"},
	Help:    "send command to remote host",
	Usage:   "shell [-h | --help] [-l | --list] <command>",
	Flags: func(f *grumble.Flags) {
		f.Bool("l", "list", false, "list command cheatsheet")
		f.Int("t", "timeout", 1, "timeout (s)")
	},
	Args: func(a *grumble.Args) {
		a.StringList("command", "command sent to remote host")
	},
	Completer: func(prefix string, args []string) []string {
		var cheatSheetCmds []string
		var yamlFile string
		if SelectVer[:7] == "windows" {
			yamlFile = "3rd_party/windows/cmd_cheatsheet.yaml"
		} else {
			yamlFile = "3rd_party/linux/cmd_cheatsheet.yaml"
		}
		yamlFile, _ = filepath.Abs(yamlFile)
		cmdCheatSheetYaml := shellopt.ReadYamlFile(yamlFile)
		for _, cmdCheatSheetStruct := range cmdCheatSheetYaml.CmdCheatSheetStructs {
			cheatSheetCmds = append(cheatSheetCmds, cmdCheatSheetStruct.Cmd)
		}
		if len(args) == 0 {
			return filterStringWithPrefix(cheatSheetCmds, prefix)
		}
		return []string{}
	},
	Run: func(c *grumble.Context) error {
		messageId := util.GenUUID()
		common.MessageQueue = append(common.MessageQueue, messageId)
		timeout := c.Flags.Int("timeout")
		listFlag := c.Flags.Bool("list")
		if listFlag {
			var yamlFile string
			if SelectVer[:7] == "windows" {
				yamlFile = "3rd_party/windows/cmd_cheatsheet.yaml"
			} else {
				yamlFile = "3rd_party/linux/cmd_cheatsheet.yaml"
			}
			yamlFile, _ = filepath.Abs(yamlFile)
			cmdCheatSheetYaml := shellopt.ReadYamlFile(yamlFile)
			shellopt.PrintTable(cmdCheatSheetYaml.CmdCheatSheetStructs)
			return nil
		}
		command := c.Args.StringList("command")
		cmdStr := strings.Join(command, " ")

		// 发送命令消息
		retData := shellopt.SendExecShellMsg(SelectClientId, cmdStr, messageId)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "shell request failed")
			return nil
		}

		// 从管道中接收消息，并打印
		select {
		case msg := <-common.ExecShellMsgChan:
			shellopt.PrintShellOutput(msg)
			common.MessageQueue = common.MessageQueue[1:]
		case <-time.After(time.Duration(10+timeout) * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}
