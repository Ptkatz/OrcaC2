package controller

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/sshopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/api"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/schollz/progressbar/v3"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var sshCmd = &grumble.Command{
	Name:  "ssh",
	Help:  "connects to target host over the SSH protocol",
	Usage: "ssh set|show|run|upload|download|tunnel [-h | --help]",
}

var sshSetCmd = &grumble.Command{
	Name:  "set",
	Help:  "set options for ssh",
	Usage: "ssh set [-h | --help] [-H | --host host] [-u | --user username] [-p | --pass password] ",
	Flags: func(f *grumble.Flags) {
		f.String("H", "host", "", "ssh host")
		f.String("u", "user", "", "ssh username")
		f.String("p", "pass", "", "ssh password")
	},
	Run: func(c *grumble.Context) error {
		host := c.Flags.String("host")
		user := c.Flags.String("user")
		pass := c.Flags.String("pass")
		if host == "" {
			host = sshopt.MySsh.SSHHost
		}
		if user == "" {
			user = sshopt.MySsh.SSHUser
		}
		if pass == "" {
			pass = sshopt.MySsh.SSHPwd
		}
		_, err := net.ResolveTCPAddr("tcp4", host)
		if err != nil {
			_, err := net.ResolveIPAddr("ip4", host)
			if err != nil {
				message := fmt.Sprintf("Invalid host address [%s]", err)
				colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
				return nil
			}
			host = fmt.Sprintf("%s:%d", host, 22)
		}
		node := SelectClientId
		if SelectId <= 0 {
			node = "Server"
		}
		sshopt.SshSet(node, host, user, pass)
		myssh, _ := json.Marshal(sshopt.MySsh)
		data, _ := crypto.Encrypt(myssh, []byte(config.AesKey))
		retData := common.SendSuccessMsg(sshopt.MySsh.Node, common.ClientId, "sshConnTest", data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, retData := common.SettleRetDataEx(msg)
			_, _, decData := common.SettleRetData(msg)
			if retData.Code != retcode.SUCCESS {
				sshopt.InitSshOption(SelectClientId)
			}
			fmt.Println(decData)
		case <-time.After(30 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var sshShowCmd = &grumble.Command{
	Name:  "show",
	Help:  "show ssh options",
	Usage: "ssh show [-h | --help]",
	Run: func(c *grumble.Context) error {
		fmt.Println("Node:", sshopt.MySsh.Node)
		fmt.Println("Target_Host:", sshopt.MySsh.SSHHost)
		fmt.Println("Username:", sshopt.MySsh.SSHUser)
		fmt.Println("Password:", sshopt.MySsh.SSHPwd)
		return nil
	},
}

var sshRunCmd = &grumble.Command{
	Name:    "run",
	Aliases: []string{"exec"},
	Help:    "connect to ssh host and execute commands",
	Usage:   "ssh run [-h | --help] <command>",
	Args: func(a *grumble.Args) {
		a.StringList("command", "command sent to remote ssh host")
	},
	Run: func(c *grumble.Context) error {
		command := c.Args.StringList("command")
		cmdStr := strings.Join(command, " ")
		sshcmd := sshopt.SshRunStruct{
			SshStruct: sshopt.MySsh,
			Command:   cmdStr,
		}
		myssh, _ := json.Marshal(sshcmd)
		data, _ := crypto.Encrypt(myssh, []byte(config.AesKey))
		retData := common.SendSuccessMsg(sshopt.MySsh.Node, common.ClientId, "sshRun", data)
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
		return nil
	},
}

var sshUploadCmd = &grumble.Command{
	Name:  "upload",
	Help:  "upload files via ssh connection",
	Usage: "file upload [-h | --help] [-l | --local local_file] [-r | remote remote_file]",
	Flags: func(f *grumble.Flags) {
		f.String("l", "local", "", "the local file path to upload")
		f.String("r", "remote", "", "uploaded to the remote file path")
	},
	Run: func(c *grumble.Context) error {
		localFile := c.Flags.String("local")
		remoteFile := c.Flags.String("remote")
		if !fileopt.IsLocalFileLegal(localFile) {
			return nil
		}
		if remoteFile == "" {
			_, fileName := filepath.Split(localFile)
			remoteFile = fileName
		}
		// 发送文件元信息
		data := sshopt.GetFileMetaInfo(localFile, remoteFile)
		retData := fileopt.SendFileMetaMsg(sshopt.MySsh.Node, data, "sshUpload")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file upload failed")
			return nil
		}
		// 分片发送文件
		fileopt.SendFileData(sshopt.MySsh.Node, localFile)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, decData := common.SettleRetData(msg)
			fmt.Println(decData)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}

var sshDownloadCmd = &grumble.Command{
	Name:  "download",
	Help:  "download files via ssh connection",
	Usage: "ssh download [-h | --help] [-r | remote remote_file] [-l | --local local_file]",
	Flags: func(f *grumble.Flags) {
		f.String("r", "remote", "", "download from remote file path")
		f.String("l", "local", "", "download to local file path")
	},
	Run: func(c *grumble.Context) error {
		remoteFile := c.Flags.String("remote")
		localFile := c.Flags.String("local")
		if !fileopt.IsRemoteFileLegal(remoteFile) {
			return nil
		}
		if localFile == "" {
			_, fileName := filepath.Split(remoteFile)
			localFile = fileName
		}
		retData := sshopt.SendDownloadRequestMsg(sshopt.MySsh.Node, common.ClientId, remoteFile, localFile, sshopt.MySsh)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "download request failed")
			return nil
		}
		var fileMetaInfo fileopt.FileMetaInfo
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, retData := common.SettleRetDataEx(msg)
			_, _, decData := common.SettleRetData(msg)
			retCode := retData.Code
			if retCode != retcode.SUCCESS {
				return nil
			}
			json.Unmarshal([]byte(decData), &fileMetaInfo)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		sliceNum := fileMetaInfo.SliceNum
		sliceSize := fileMetaInfo.SliceSize
		remainSize := fileMetaInfo.RemainSize
		pSaveFile, _ := os.OpenFile(localFile, os.O_CREATE|os.O_RDWR, 0600)
		defer pSaveFile.Close()
		// 循环从管道中获取文件元数据并写入
		bar := progressbar.Default(int64(sliceNum) + 1)
		for i := 0; i < sliceNum+1; i++ {
			var err error
			select {
			case metaData := <-common.FileSliceMsgChan:
				_, err = pSaveFile.Write(metaData)
			case <-time.After(8 * time.Second):
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "file download failed")
				return nil
			}
			bar.Add(1)
			if err != nil {
				break
			}
		}
		fileSize := fileopt.GetFileSize(localFile)
		if fileSize == int64(sliceNum)*sliceSize+remainSize {
			colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "file download success")
		} else {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file download failed")
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}

var sshTunnelCmd = &grumble.Command{
	Name:  "tunnel",
	Help:  "ssh tunnel",
	Usage: "ssh tunnel [-h | --help]",
}

var sshTunnelStartCmd = &grumble.Command{
	Name: "start",
	Help: "ssh tunnel start",
	Usage: "ssh tunnel start [-h | --help] [-s | --source source_host] [-t | --target target_host] \n" +
		"  eg: \n   ssh tunnel start -t 192.168.1.10:3306 -s 10.1.1.1:13306",
	Flags: func(f *grumble.Flags) {
		f.String("s", "source", "", "source host: Server address&port mapped to")
		f.String("t", "target", "", "target host: The mapped target address&port")
	},
	Run: func(c *grumble.Context) error {
		if sshopt.MySsh.SSHHost == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please set up ssh first")
			return nil
		}
		source := c.Flags.String("source")
		target := c.Flags.String("target")
		if source == "" {
			if SelectId <= 0 {
				ip, _, _ := strings.Cut(api.HOST, ":")
				source = fmt.Sprintf("%s:%d", ip, 8443)
			} else {
				source = fmt.Sprintf("%s:%d", SelectIp, 8443)
			}
		}
		if target == "" {
			target = sshopt.MySsh.SSHHost
		}
		var err error
		_, err = net.ResolveTCPAddr("tcp4", source)
		if err != nil {
			message := fmt.Sprintf("Invalid source host address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		_, err = net.ResolveTCPAddr("tcp4", target)
		if err != nil {
			message := fmt.Sprintf("Invalid target host address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		sshTunnelStruct := sshopt.SshTunnelStruct{
			SshStruct: sshopt.MySsh,
			Source:    source,
			Target:    target,
		}
		marshal, _ := json.Marshal(sshTunnelStruct)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(sshopt.MySsh.Node, common.ClientId, "sshTunnelStart", data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, decData := common.SettleRetData(msg)
			fmt.Println(decData)
		case <-time.After(30 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}

var sshTunnelListCmd = &grumble.Command{
	Name:  "list",
	Help:  "list ssh tunnel",
	Usage: "ssh tunnel list [-h | --help]",
	Run: func(c *grumble.Context) error {
		retData := common.SendSuccessMsg("Server", common.ClientId, "sshTunnelList", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			sshTunnelLists := sshopt.GetSshTunnelLists(msg)
			sshopt.PrintSshTunnelTable(sshTunnelLists)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		return nil
	},
}

var sshTunnelCloseCmd = &grumble.Command{
	Name:  "close",
	Help:  "turn off the ssh tunnel by id",
	Usage: "ssh tunnel close [-h | --help] <id>",
	Args: func(a *grumble.Args) {
		a.Int("id", "close ssh tunnel id")
	},
	Run: func(c *grumble.Context) error {
		tunnelId := c.Args.Int("id")
		var uid string
		retData := common.SendSuccessMsg("Server", common.ClientId, "sshTunnelList", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		var sshTunnelList []sshopt.SshtunnelList
		select {
		case msg := <-common.DefaultMsgChan:
			data := common.GetHttpRetData(msg)
			json.Unmarshal([]byte(data), &sshTunnelList)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		closeFlag := false
		for i, record := range sshTunnelList {
			if tunnelId == i+1 {
				uid = record.Uid
				closeFlag = true
			}
		}
		if !closeFlag {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "The corresponding id does not exist")
			return nil
		}
		data, _ := crypto.Encrypt([]byte(uid), []byte(config.AesKey))
		retData = common.SendSuccessMsg(sshopt.MySsh.Node, common.ClientId, "sshTunnelClose", data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "ssh tunnel closed successfully")
		return nil
	},
}
