package infoopt

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

// 获取杀软信息
//func GetAvName() string {
//	ac, err := shellopt.Cmd("wmic /namespace:\\\\root\\securitycenter2 path antivirusproduct get displayname /Format:List")
//	if err != nil {
//		return ""
//	}
//	b := strings.Replace(ac, "\n", "", -1)
//	c := strings.Replace(b, "\r\r\r\r\r\r", "\n", -1)
//	d := strings.Replace(c, "\r", "", -1)
//	AvName := strings.Replace(d, "displayName=", " ", -1)
//	return AvName
//}

// 网卡配置信息
func ifconfig() (stdout string, err error) {
	fSize := uint32(0)
	b := make([]byte, 1000)

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var adapterInfo *syscall.IpAdapterInfo
	adapterInfo = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	err = syscall.GetAdaptersInfo(adapterInfo, &fSize)

	// Call it once to see how much data you need in fSize
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b := make([]byte, fSize)
		adapterInfo = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(adapterInfo, &fSize)
		if err != nil {
			return "", err
		}
	}

	for _, iface := range ifaces {
		for ainfo := adapterInfo; ainfo != nil; ainfo = ainfo.Next {
			if int(ainfo.Index) == iface.Index {
				stdout += fmt.Sprintf("%-30s", fmt.Sprintf("Name:%s", iface.Name))
				stdout += fmt.Sprintf("%-30s", fmt.Sprintf("MAC-Address:%s", iface.HardwareAddr.String()))
				ipentry := &ainfo.IpAddressList
				for ; ipentry != nil; ipentry = ipentry.Next {
					stdout += fmt.Sprintf("%-30s", fmt.Sprintf("\nIP-Address:%s\n", ipentry.IpAddress.String))
					stdout += fmt.Sprintf("%-30s", fmt.Sprintf("\nSubnet-Mask:%s\n", ipentry.IpMask.String))
				}
				gateways := &ainfo.GatewayList
				for ; gateways != nil; gateways = gateways.Next {
					stdout += fmt.Sprintf("%-30s", fmt.Sprintf("\nGateway:%s\n", gateways.IpAddress.String))
				}

				if ainfo.DhcpEnabled != 0 {
					stdout += fmt.Sprintf("\nDHCP--Enabled\n")
					dhcpServers := &ainfo.DhcpServer
					for ; dhcpServers != nil; dhcpServers = dhcpServers.Next {
						stdout += fmt.Sprintf("%-30s", fmt.Sprintf("\nDHCP-Server:%s\n", dhcpServers.IpAddress.String))
					}
				} else {
					stdout += fmt.Sprintf("%-30s", fmt.Sprintf("\nDHCP--Disabled\n"))
				}
				stdout += "\n"
			}
		}
	}

	return stdout, nil
}
