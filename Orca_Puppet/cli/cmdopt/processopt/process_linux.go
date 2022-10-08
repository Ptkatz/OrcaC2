package processopt

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"io"
	"os"
)

func GetProcs() (procs []ProcessInfo, err error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range processes {
		pid := p.Pid
		var proc = ProcessInfo{
			Name:     getProcessName(uint32(pid)),
			Pid:      uint32(pid),
			Ppid:     GetProcessParrentId(uint32(pid)),
			Mem:      getProcessMemUsage(uint32(pid)),
			Username: getProcessUsername(uint32(pid)),
			Arch:     getProcessArchitecture(uint32(pid)),
		}
		procs = append(procs, proc)
	}
	return
}

func getProcessArchitecture(pid uint32) string {
	exePath := fmt.Sprintf("/proc/%d/exe", pid)

	f, err := os.Open(exePath)
	if err != nil {
		return "*"
	}
	_, err = f.Seek(0x12, 0)
	if err != nil {
		return "*"
	}
	mach := make([]byte, 2)
	n, err := io.ReadAtLeast(f, mach, 2)

	f.Close()

	if err != nil || n < 2 {
		return "*"
	}

	if mach[0] == 0xb3 {
		return "aarch64"
	}
	if mach[0] == 0x03 {
		return "x86"
	}
	if mach[0] == 0x3e {
		return "x86_64"
	}
	return "*"
}
