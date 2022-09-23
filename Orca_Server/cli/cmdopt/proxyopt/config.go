package proxyopt

import "github.com/esrrhs/go-engine/src/proxy"

type ProxyServerParam struct {
	Proto  string
	Listen string
	Key    string
}

type ProxyServer struct {
	Uid    string
	Server *proxy.Server
	Param  ProxyServerParam
}

type ProxyClientParam struct {
	Proto      string
	ProxyProto string
	ConnType   string
	ServerAddr string
	Key        string
	FromAddr   string
	ToAddr     string
}

type ProxyClient struct {
	Uid              string
	Client           *string
	ClientId         string
	ProxyClientParam ProxyClientParam
}

var ProxyServerLists = make([]ProxyServer, 0)
