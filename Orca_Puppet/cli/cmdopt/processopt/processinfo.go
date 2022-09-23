package processopt

import "github.com/shirou/gopsutil/v3/process"

type ProcessInfo struct {
	Name     string
	Pid      uint32
	Ppid     uint32
	Mem      uint64
	Username string
	Arch     string
}

// 获取进程占用内存
func getProcessMemUsage(pid uint32) uint64 {
	newProcess, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0
	}
	meminfo, err := newProcess.MemoryInfo()
	if err != nil {
		return 0
	}
	return meminfo.RSS / 1024
}

// 获取进程用户名
func getProcessUsername(pid uint32) string {
	newProcess, err := process.NewProcess(int32(pid))
	if err != nil {
		return ""
	}
	username, err := newProcess.Username()
	if err != nil {
		return ""
	}
	return username
}

// 获取进程名
func getProcessName(pid uint32) string {
	newProcess, err := process.NewProcess(int32(pid))
	if err != nil {
		return ""
	}
	name, err := newProcess.Name()
	if err != nil {
		return ""
	}
	return name
}

// 获取父进程id
func GetProcessParrentId(pid uint32) uint32 {
	newProcess, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0
	}
	ppid, err := newProcess.Ppid()
	if err != nil {
		return 0
	}
	return uint32(ppid)
}
