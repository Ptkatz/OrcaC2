package sshopt

import (
	"Orca_Master/cli/common"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type SshTunnelStruct struct {
	SshStruct SshOption
	Source    string
	Target    string
}

type SshtunnelList struct {
	gorm.Model
	Uid      string `json:"uid" gorm:"unique"`
	ClientId string `json:"clientId"`
	Source   string `json:"source"`
	Target   string `json:"target"`
}

func PrintSshTunnelTable(sshTunnelLists []SshtunnelList) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(48)

	table.SetHeader([]string{"id", "Uid", "clientId", "source", "target"})
	for i, sshTunnelList := range sshTunnelLists {
		data = append(data, []string{strconv.Itoa(i + 1), sshTunnelList.Uid, sshTunnelList.ClientId, sshTunnelList.Source, sshTunnelList.Target})
	}
	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}

func GetSshTunnelLists(message string) []SshtunnelList {
	var retData common.RetData
	var sshTunnelLists []SshtunnelList
	err := json.Unmarshal([]byte(message), &retData)
	data := retData.Data.(string)
	err = json.Unmarshal([]byte(data), &sshTunnelLists)
	if err != nil {
		return nil
	}
	return sshTunnelLists
}
