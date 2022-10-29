package listopt

import (
	"Orca_Master/cli/common"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

type HostList struct {
	gorm.Model
	ClientId  string `json:"clientId" gorm:"unique"` //唯一标识
	Hostname  string `json:"hostname"`               //主机名
	Ip        string `json:"ip"`                     //上线ip
	ConnPort  string `json:"connPort"`               // 上线端口
	Privilege string `json:"privilege"`              //执行权限
	Os        string `json:"os"`                     //系统版本
	Version   string `json:"version"`                //上线客户端版本
	Remarks   string `json:"remarks"`                //备注
}

// 从消息字符串中获取HostList切片数组
func GetHostLists(message string) []HostList {
	var retData common.RetData
	var hostlists []HostList
	err := json.Unmarshal([]byte(message), &retData)
	data := retData.Data.(string)
	err = json.Unmarshal([]byte(data), &hostlists)
	if err != nil {
		return nil
	}
	return hostlists
}

// 打印表格
func PrintTable(hostlists []HostList, identity int) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(48)

	table.SetHeader([]string{"id", "Hostname", "Ip", "Os", "Arch", "Privilege", "Port"})
	for i, host := range hostlists {
		verSplit := strings.Split(host.Version, ":")
		data = append(data, []string{strconv.Itoa(i + 1), host.Hostname, host.Ip, host.Os, verSplit[2], host.Privilege, host.ConnPort})
	}

	for _, raw := range data {
		if identity != -1 {
			if raw[0] == strconv.Itoa(identity) {
				table.Append(raw)
				break
			}
		} else {
			table.Append(raw)
		}
	}
	table.Render()
	return
}
