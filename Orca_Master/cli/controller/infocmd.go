package controller

import (
	"Orca_Master/cli/cmdopt/infoopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/retcode"
	"encoding/json"
	"github.com/desertbit/grumble"
	"time"
)

var infoCmd = &grumble.Command{
	Name:  "info",
	Help:  "get basic information of remote host",
	Usage: "info [-h | --help]",
	Run: func(c *grumble.Context) error {
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "info", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "get info request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			var info infoopt.Info
			json.Unmarshal([]byte(msg), &info)
			infoopt.PrintClientInfoTable(info)
			infoopt.PrintSystemInfoTable(info)
		case <-time.After(20 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}
