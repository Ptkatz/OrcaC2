package main

import (
	"Orca_Master/cli/controller"
	"Orca_Master/define/api"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/login"
	"Orca_Master/stager"
	"fmt"
	"github.com/desertbit/grumble"
	"net"
	"os"
)

func main() {
	// 进入命令行终端前的准备操作
	controller.App.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		host := flags.String("host")
		username := flags.String("username")
		password := flags.String("password")
		// 判断服务端ip与端口是否输入正确
		_, err := net.ResolveTCPAddr("tcp4", host)
		if err != nil {
			message := fmt.Sprintf("Invalid OrcaServer TCP address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			os.Exit(1)
		}
		api.InitApi(host)
		retData := login.UserLogin(api.MASTER_LOGIN_API, username, password)
		if retData.Code == retcode.SUCCESS {
			controller.Uname = username
			colorcode.PrintMessage(colorcode.SIGN_NOTICE, retData.Msg)
			config.AesKey = retData.Data.(string)
			controller.InitPrompt = fmt.Sprintf("Orca[%s] » ", controller.Uname)
			controller.App.SetPrompt(controller.InitPrompt)
			go stager.Init()
		} else {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, retData.Msg)
			os.Exit(1)
		}
		return nil
	})
	grumble.Main(controller.App)
}
