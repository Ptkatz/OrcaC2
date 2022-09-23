package processopt

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const TH32CS_SNAPPROCESS = 0x00000002

// 获取进程列表
func GetProcs() (procs []ProcessInfo, err error) {
	snap := createToolhelp32Snapshot(TH32CS_SNAPPROCESS, uint32(0))
	if snap == 0 {
		err = syscall.GetLastError()
		return
	}

	defer windows.CloseHandle(snap)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))
	if nil != windows.Process32First(snap, &pe32) {
		err = syscall.GetLastError()
		return
	}
	arch := detArch(pe32.ProcessID)
	procs = append(procs, ProcessInfo{syscall.UTF16ToString(pe32.ExeFile[:260]), pe32.ProcessID, pe32.ParentProcessID, getProcessMemUsage(pe32.ProcessID), getProcessUsername(pe32.ProcessID), arch})
	for process32Next(snap, &pe32) {
		arch = detArch(pe32.ProcessID)
		procs = append(procs, ProcessInfo{syscall.UTF16ToString(pe32.ExeFile[:260]), pe32.ProcessID, pe32.ParentProcessID, getProcessMemUsage(pe32.ProcessID), getProcessUsername(pe32.ProcessID), arch})
	}
	return
}

// 创建进程快照
func createToolhelp32Snapshot(flags, processId uint32) windows.Handle {
	ret, _ := windows.CreateToolhelp32Snapshot(
		flags,
		processId)
	if ret <= 0 {
		return windows.Handle(0)
	}
	return windows.Handle(ret)
}

// 查找下一进程信息
func process32Next(snapshot windows.Handle, pe *windows.ProcessEntry32) bool {
	err := windows.Process32Next(
		snapshot,
		pe)
	return nil == err
}

// 确定程序架构
func detArch(pid uint32) string {
	arch := "x86"
	processHandle, err := syscall.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return "*"
	}
	var wow64Process bool
	kernel32 := windows.NewLazySystemDLL("kernel32")
	procIsWow64Process := kernel32.NewProc("IsWow64Process")

	r1, _, _ := procIsWow64Process.Call(
		uintptr(processHandle),
		uintptr(unsafe.Pointer(&wow64Process)))
	if int(r1) == 0 {
		return "*"
	}
	if !wow64Process {
		arch = "x86_64"
	}
	return arch
}
