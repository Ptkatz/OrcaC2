package proxyopt

import (
	"Orca_Master/cli/common"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type ProxyClientParam struct {
	Proto      string `json:"proto"`
	ProxyProto string `json:"proxyProto"`
	ConnType   string `json:"connType"`
	ServerAddr string `json:"serverAddr"`
	Key        string `json:"key"`
	FromAddr   string `json:"fromAddr"`
	ToAddr     string `json:"toAddr"`
}

type ProxyClientList struct {
	gorm.Model
	Uid        string `json:"uid" gorm:"unique"`
	ClientId   string `json:"clientId"`
	Proto      string `json:"proto"`
	ProxyProto string `json:"proxyProto"`
	ConnType   string `json:"connType"`
	ServerAddr string `json:"serverAddr"`
	Key        string `json:"key"`
	FromAddr   string `json:"fromAddr"`
	ToAddr     string `json:"toAddr"`
}

var ReliableConnType = []string{
	"bind",
	"reverse",
	"bind_socks5",
	"reverse_socks5",
}

func HasReliableConnType(dst string) bool {
	for _, i := range ReliableConnType {
		if i == dst {
			return true
		}
	}
	return false
}

func SupportReliableConnType() []string {
	ret := make([]string, 0)
	ret = append(ret, "bind")
	ret = append(ret, "reverse")
	ret = append(ret, "bind_socks5")
	ret = append(ret, "reverse_socks5")
	return ret
}

func PrintClientTable(proxyLists []ProxyClientList) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(48)

	table.SetHeader([]string{"id", "Uid", "clientId", "proto", "proxyProto", "connType", "serverAddr", "fromAddr", "toAddr"})
	for i, proxyList := range proxyLists {
		data = append(data, []string{strconv.Itoa(i + 1), proxyList.Uid, proxyList.ClientId, proxyList.Proto, proxyList.ProxyProto, proxyList.ConnType, proxyList.ServerAddr, proxyList.FromAddr, proxyList.ToAddr})
	}
	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}

func GetClientProxyLists(message string) []ProxyClientList {
	var retData common.RetData
	var proxyLists []ProxyClientList
	err := json.Unmarshal([]byte(message), &retData)
	data := retData.Data.(string)
	err = json.Unmarshal([]byte(data), &proxyLists)
	if err != nil {
		return nil
	}
	return proxyLists
}
