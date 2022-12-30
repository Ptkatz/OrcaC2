package api

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	IP                 = "4294967295"
	PORT               = ":65535"
	HOST               string
	CONN_SERVER_API    string
	REGISTER_API       string
	SEND_TO_CLIENT_API string
)

func InitApi(host string) {
	ip, port, _ := strings.Cut(host, ":")
	IP = strconv.FormatInt(IpAddrToInt(ip), 10)
	nPort, _ := strconv.Atoi(port)
	PORT = fmt.Sprintf(":%05d", nPort)
	HOST = ParseHost()
	CONN_SERVER_API = "ws://" + HOST + "/ws?systemId="
	REGISTER_API = "http://" + HOST + "/api/register"
	SEND_TO_CLIENT_API = "http://" + HOST + "/api/send_to_client"
}

func ParseHost() string {
	i, _ := strconv.ParseUint(IP, 10, 64)
	xIp := IntToIpAddr(int64(i))
	xPort, _ := strconv.ParseUint(PORT[1:], 10, 32)
	return fmt.Sprintf("%s:%d", xIp, xPort)
}

func IpAddrToInt(ipAddr string) int64 {
	bits := strings.Split(ipAddr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

func IntToIpAddr(intIP int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
