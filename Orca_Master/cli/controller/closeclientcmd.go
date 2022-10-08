package controller

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/retcode"
	"fmt"
	"github.com/desertbit/grumble"
	"time"
)

var closeClientCmd = &grumble.Command{
	Name:  "close",
	Help:  "close the selected remote client",
	Usage: "close [-h | --help]",
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		currentTime := time.Now().Format("2006/01/02 15:04:05")
		timeSign := colorcode.COLOR_GREY + currentTime + colorcode.END
		fmt.Printf("%s %s do you want to close the remote client? [y/N]: ", timeSign, colorcode.SIGN_QUEST)
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "yes" {
			retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "closeClient", "")
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "close client request failed")
				return nil
			}
			colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "the remote client has been closed, waiting for the server to respond")
			BackMainMenu()
			return nil
		}
		return nil
	},
}
