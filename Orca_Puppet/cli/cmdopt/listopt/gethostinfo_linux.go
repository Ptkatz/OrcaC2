package listopt

import (
	"Orca_Puppet/cli/cmdopt/shellopt"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

var sysType string = runtime.GOOS

// 获取操作系统信息
func GetOsName() string {
	ac, err := shellopt.Cmd("cat /etc/os-release | grep PRETTY_NAME | awk -F \"[\\\"\\\"]\" '{print $2}'")
	if err != nil {
		return ""
	}
	sysVersion := strings.Replace(ac, "\n", "", -1)
	return sysVersion
}

// 获取当前执行身份 (user/admin/system)
func GetExecPrivilege() string {
	status := "user"
	command := "id -u"
	ac, _ := shellopt.Cmd(command)
	id := strings.Replace(ac, "\n", "", -1)
	if id == "0" {
		status = "root"
	}
	return status
}

func GetHostName() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	hostname, _ := os.Hostname()
	hostname += "/" + u.Username
	return hostname
}

func GetConnPort() (string, error) {
	pid := os.Getpid()
	cmd := fmt.Sprintf("netstat -anopt |grep %d|head -n 1|awk '{printf $4}'|awk -F '[:]' '{print $NF}'", pid)
	ac, err := shellopt.Cmd(cmd)
	if err != nil {
		return "", err
	}
	result := strings.Replace(ac, "\n", "", -1)
	re := regexp.MustCompile("[0-9]+")
	ports := re.FindAllString(result, -1)
	return ports[0], nil
}
