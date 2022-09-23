package sqlmgmt

import (
	"Orca_Server/cli/cmdopt/sshopt"
	"encoding/json"
	"gorm.io/gorm"
)

func InitSshTunnel() {
	Db = GetDb()
	Db.Exec("DELETE FROM sshtunnel_lists")
}

func DelSshTunnelRecordByUid(uid string) {
	Db = GetDb()
	Db.Exec("DELETE FROM sshtunnel_lists Where uid = ?", uid)
}

func DelSshTunnelRecordByClientId(clientID string) {
	Db = GetDb()
	Db.Exec("DELETE FROM sshtunnel_lists Where client_id = ?", clientID)
}

func AddSshTunnel(client sshopt.SshTunnelBaseRecord) {
	Db = GetDb()
	uid := client.Uid
	clientId := client.ClientId
	source := client.Source
	target := client.Target

	sshTunnelList := SshtunnelList{
		Model:    gorm.Model{},
		Uid:      uid,
		ClientId: clientId,
		Source:   source,
		Target:   target,
	}
	Db.Create(&sshTunnelList)
}

func ListSshTunnel() []byte {
	Db = GetDb()
	var sshTunnelLists []SshtunnelList
	Db.Find(&sshTunnelLists)
	listMsg, _ := json.Marshal(sshTunnelLists)
	return listMsg
}
