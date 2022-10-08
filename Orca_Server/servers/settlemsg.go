package servers

import (
	"Orca_Server/cli/common/setchannel"
	"Orca_Server/setting"
	"Orca_Server/tools/crypto"
)

func SettleMsg(msg, clientId string, data *string) {
	decData, _ := crypto.Decrypt(*data, []byte(setting.CommonSetting.CryptoKey))
	switch msg {
	case "hostInfo":
		go hostCmd(decData)
		return
	case "listHosts":
		go listCmd(clientId)
		return
	case "ipSearch":
		go ipSearchCmd(*data, clientId)
		return
	case "proxyServerStart":
		go proxyServerStartCmd(decData, clientId)
		return
	case "proxyServerList":
		go proxyServerListCmd(clientId)
		return
	case "proxyServerClose":
		go proxyServerCloseCmd(decData)
		return
	case "proxyClientStart":
		go proxyClientStartCmd(decData, clientId)
		return
	case "proxyClientList":
		go proxyClientListCmd(clientId)
		return
	case "proxyClientClose":
		go proxyClientCloseCmd(decData)
		return
	case "sshConnTest":
		go sshConnTestCmd(clientId, decData)
		return
	case "sshRun":
		go sshRunCmd(clientId, decData)
		return
	case "sshUpload":
		go sshUploadCmd(clientId, decData)
		return
	case "sshDownload":
		go sshDownloadCmd(clientId, decData)
		return
	case "sshTunnelStart":
		go sshTunnelStartCmd(clientId, decData)
		return
	case "sshTunnelList":
		go sshTunnelListCmd(clientId)
		return
	case "sshTunnelClose":
		go sshTunnelCloseCmd(decData)
		return
	case "sshTunnelDel":
		go sshTunnelDelCmd(decData)
		return
	case "sshTunnelAdd":
		go sshTunnelAddCmd(decData)
		return
	case "portScan":
		go portScanCmd(clientId, decData)
		return
	case "portCrack":
		go portCrackCmd(clientId, decData)
		return
	case "fileUpload":
		go fileUploadCmd(clientId, decData)
		return

	case "sliceData":
		m, exist := setchannel.GetFileSliceDataChan(clientId)
		if !exist {
			m = make(chan interface{})
			setchannel.AddFileSliceDataChan(clientId, m)
		}
		m <- []byte(decData)
		return
	}
}
