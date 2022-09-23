package proxyopt

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

type ProxyServerParam struct {
	Proto  string
	Listen string
	Key    string
}

var ReliableProto = []string{
	"tcp",
	"rudp",
	"ricmp",
	"kcp",
	"quic",
	"rhttp",
}

type ProxyServer struct {
	Uid    string
	Server *int
	Param  ProxyServerParam
}

func HasReliableProto(dst string) bool {
	for _, i := range ReliableProto {
		if i == dst {
			return true
		}
	}
	return false
}

func SupportReliableProtos() []string {
	ret := make([]string, 0)
	ret = append(ret, "tcp")
	ret = append(ret, "rudp")
	ret = append(ret, "ricmp")
	ret = append(ret, "kcp")
	ret = append(ret, "quic")
	ret = append(ret, "rhttp")
	return ret
}

func PrintServerTable(proxyList []ProxyServer) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(48)

	table.SetHeader([]string{"id", "Uid", "Proto", "Listen", "Key"})
	for i, proxyServer := range proxyList {
		data = append(data, []string{strconv.Itoa(i + 1), proxyServer.Uid, proxyServer.Param.Proto, proxyServer.Param.Listen, proxyServer.Param.Key})
	}
	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}
