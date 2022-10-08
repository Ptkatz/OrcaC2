package sqlmgmt

import (
	"Orca_Server/cli/cmdopt/proxyopt"
	"encoding/json"
	"gorm.io/gorm"
)

func InitProxy() {
	Db = GetDb()
	Db.Exec("DELETE FROM proxy_lists")
}

func DelProxyRecordByUid(uid string) {
	Db = GetDb()
	Db.Exec("DELETE FROM proxy_lists Where uid = ?", uid)
}

func DelProxyRecordByClientId(clientId string) {
	Db = GetDb()
	Db.Exec("DELETE FROM proxy_lists Where client_id = ?", clientId)
}

func AddProxyClient(client proxyopt.ProxyClient) {
	Db = GetDb()
	uid := client.Uid
	clientId := client.ClientId
	serverAddr := client.ProxyClientParam.ServerAddr
	proto := client.ProxyClientParam.Proto
	proxyProto := client.ProxyClientParam.ProxyProto
	connType := client.ProxyClientParam.ConnType
	fromAddr := client.ProxyClientParam.FromAddr
	toAddr := client.ProxyClientParam.ToAddr
	key := client.ProxyClientParam.Key

	proxyList := ProxyList{
		Model:      gorm.Model{},
		Uid:        uid,
		ClientId:   clientId,
		Proto:      proto,
		ProxyProto: proxyProto,
		ConnType:   connType,
		ServerAddr: serverAddr,
		Key:        key,
		FromAddr:   fromAddr,
		ToAddr:     toAddr,
	}
	Db.Create(&proxyList)
}

func ListProxy() []byte {
	Db = GetDb()
	var proxyLists []ProxyList
	Db.Find(&proxyLists)
	listMsg, _ := json.Marshal(proxyLists)
	return listMsg
}
