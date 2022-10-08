package infoopt

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
	"time"
)

type Info struct {
	ClientInfo ClientInfo
	SystemInfo SystemInfo
}

type ClientInfo struct {
	ExecPath   string
	ClientId   string
	Uptime     string
	Privilege  string
	Version    string
	Hostname   string
	CurrentPid string
	OnlineIp   string
	ConnPort   string
	ExternalIp string
	Address    string
}

type SystemInfo struct {
	HostInfo string
	CpuInfo  string
	MemInfo  string
	DiskInfo string
	IfConfig string
}

func PrintClientInfoTable(info Info) {
	fmt.Println("Client Info:")
	clientInfo := info.ClientInfo
	common.SendSuccessMsg("Server", common.ClientId, "ipSearch", clientInfo.ExternalIp)
	select {
	case msg := <-common.DefaultMsgChan:
		_, _, addr := common.SettleRetDataNotDec(msg)
		clientInfo.Address = addr
	case <-time.After(10 * time.Second):
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
		return
	}
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColMinWidth(0, 20)
	table.SetColMinWidth(1, 65)
	t := reflect.TypeOf(clientInfo)
	v := reflect.ValueOf(clientInfo)
	for i := 0; i < t.NumField(); i++ {
		data = append(data, []string{t.Field(i).Name, v.Field(i).String()})
	}

	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}

func PrintSystemInfoTable(info Info) {
	fmt.Println("System Info:")
	systemInfo := info.SystemInfo
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColMinWidth(0, 20)
	table.SetColMinWidth(1, 65)
	t := reflect.TypeOf(systemInfo)
	v := reflect.ValueOf(systemInfo)
	for i := 0; i < t.NumField(); i++ {
		data = append(data, []string{t.Field(i).Name, v.Field(i).String()})
	}
	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}
