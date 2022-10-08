package proxyopt

import "Orca_Puppet/pkg/go-engine/proxy"

type ProxyClientParam struct {
	Proto      string
	ProxyProto string
	ConnType   string
	ServerAddr string
	Key        string
	FromAddr   string
	ToAddr     string
}

var ConnTypeMap = map[string]string{
	"bind":           "PROXY",
	"reverse":        "REVERSE_PROXY",
	"bind_socks5":    "SOCKS5",
	"reverse_socks5": "REVERSE_SOCKS5",
}

type ProxyClient struct {
	Uid              string
	Client           *proxy.Client
	ClientId         string
	ProxyClientParam ProxyClientParam
}

var ProxyClientList = make([]ProxyClient, 0)
