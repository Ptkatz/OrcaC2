package processopt

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

type ProcessInfo struct {
	Name     string
	Pid      uint32
	Ppid     uint32
	Mem      uint64
	Username string
	Arch     string
}

func PrintTable(processInfos []ProcessInfo, name string, pid int) {
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "pid", "ppid", "mem", "owner", "arch"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetColWidth(48)
	if name == "" && pid == -1 {
		for i, processInfo := range processInfos {
			data = append(data, []string{strconv.Itoa(i + 1), processInfo.Name, strconv.Itoa(int(processInfo.Pid)), strconv.Itoa(int(processInfo.Ppid)), strconv.Itoa(int(processInfo.Mem)) + " KB", processInfo.Username, processInfo.Arch})
		}
	} else if name != "" && pid != -1 {
		i := 0
		for _, processInfo := range processInfos {
			if processInfo.Name == name && int(processInfo.Pid) == pid {
				data = append(data, []string{strconv.Itoa(i + 1), processInfo.Name, strconv.Itoa(int(processInfo.Pid)), strconv.Itoa(int(processInfo.Ppid)), strconv.Itoa(int(processInfo.Mem)) + " KB", processInfo.Username, processInfo.Arch})
				i++
			}
		}
	} else {
		i := 0
		for _, processInfo := range processInfos {
			if processInfo.Name == name || int(processInfo.Pid) == pid || int(processInfo.Ppid) == pid {
				data = append(data, []string{strconv.Itoa(i + 1), processInfo.Name, strconv.Itoa(int(processInfo.Pid)), strconv.Itoa(int(processInfo.Ppid)), strconv.Itoa(int(processInfo.Mem)) + " KB", processInfo.Username, processInfo.Arch})
				i++
			}
		}
	}

	for _, raw := range data {
		table.Append(raw)
	}
	table.Render()
}

// 获取进程Pid
func GetPid(name string, processInfos []ProcessInfo) uint32 {
	var err error
	if nil != err {
		return 0
	}
	for _, value := range processInfos {
		if name == value.Name {
			return value.Pid
		}
	}
	return 0
}
