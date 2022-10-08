package infoopt

import (
	"fmt"
	"net"
)

// 网卡配置信息
func ifconfig() (stdout string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		stdout += fmt.Sprintf("%-30s", fmt.Sprintf("Name:%s", i.Name))
		stdout += fmt.Sprintf("%-30s", fmt.Sprintf("MAC_Address:%s", i.HardwareAddr.String()))
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, a := range addrs {
			stdout += fmt.Sprintf("%-30s", fmt.Sprintf("IP_Address:%s", a.String()))
		}
	}
	return stdout, nil
}
