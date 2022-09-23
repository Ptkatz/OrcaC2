package sqlmgmt

import (
	"Orca_Server/cli/cmdopt/hostopt"
	"encoding/json"
	"gorm.io/gorm"
)

func InitHost() {
	Db = GetDb()
	Db.Exec("DELETE FROM host_lists")
}

func AddHost(hostInfo hostopt.HostInfo) {
	Db = GetDb()
	hostList := HostList{
		Model:     gorm.Model{},
		ClientId:  hostInfo.ClientId,
		Hostname:  hostInfo.Hostname,
		Ip:        hostInfo.Ip,
		ConnPort:  hostInfo.ConnPort,
		Privilege: hostInfo.Privilege,
		Os:        hostInfo.Os,
		Version:   hostInfo.Version,
		Remarks:   "",
	}
	Db.Create(&hostList)
}

func DelHostRecordByClientId(clientId string) {
	Db = GetDb()
	Db.Exec("DELETE FROM host_lists where client_id = ?", clientId)
}

func ListHosts() []byte {
	Db = GetDb()
	var hosts []HostList
	Db.Find(&hosts)
	listMsg, _ := json.Marshal(hosts)
	return listMsg
}
