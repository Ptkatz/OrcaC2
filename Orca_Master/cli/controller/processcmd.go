package controller

import (
	"Orca_Master/cli/cmdopt/processopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"strconv"
	"strings"
	"time"
)

var processCmd = &grumble.Command{
	Name:    "process",
	Aliases: []string{"ps"},
	Help:    "manage remote host processes",
	Usage:   "process [-h | --help] list|kill",
}

var processListCmd = &grumble.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Help:    "show process list",
	Usage:   "process list [-h | --help] [-n | --name search_by_name] [-i | --pid search_by_pid]",
	Flags: func(f *grumble.Flags) {
		f.String("n", "name", "", "process name for search")
		f.Int("i", "pid", -1, "pid or ppid for search")
	},
	Run: func(c *grumble.Context) error {
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "processList", "", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "message request failed")
			return nil
		}
		pid := c.Flags.Int("pid")
		name := strings.Trim(c.Flags.String("name"), " ")

		if SelectVer[:7] == "windows" {
			if name != "" {
				if len(name) > 4 {
					if name[len(name)-4:len(name)] != ".exe" {
						name = name + ".exe"
					}
				} else {
					name = name + ".exe"
				}
			}
		}
		var processInfos []processopt.ProcessInfo
		select {
		case msg := <-common.ProcessListChan:
			json.Unmarshal([]byte(msg), &processInfos)
		case <-time.After(5 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		processopt.PrintTable(processInfos, name, pid)
		return nil
	},
}

var processKillCmd = &grumble.Command{
	Name: "kill",
	Help: "terminate the specified process",
	Args: func(a *grumble.Args) {
		a.Int("pid", "process id to terminate")
	},
	Run: func(c *grumble.Context) error {
		pid := c.Args.Int("pid")
		data, _ := crypto.Encrypt([]byte(strconv.Itoa(pid)), []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "processKill", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "message request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			fmt.Println(msg)
		case <-time.After(5 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}
