package controller

import (
	"Orca_Master/cli/cmdopt/proxyopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
	"net"
	"time"
)

var proxyCmd = &grumble.Command{
	Name: "proxy",
	Help: "activate the proxy function",
	Usage: "proxy [-h | --help] server|client\n\n" +
		"  * Supported protocol: TCP, UDP, RUDP (Reliable UDP), RICMP (Reliable ICMP), RHTTP (Reliable HTTP), KCP, Quic\n  * Support type: forward proxy, reverse agent, SOCKS5 forward agent, SOCKS5 reverse agent",
}

var proxyServerCmd = &grumble.Command{
	Name:  "server",
	Help:  "activate the proxy server function",
	Usage: "proxy server [-h | --help]",
}

var proxyServerStartCmd = &grumble.Command{
	Name: "start",
	Help: "the server starts the proxy listener",
	Usage: "proxy server start [-h | --help] [-p | --proto proto] [-l | --listen listen] [-k | --key key]\n" +
		"  eg: \n   proxy server start -l 8808 -k admin@123",
	Flags: func(f *grumble.Flags) {
		f.String("p", "proto", "tcp", "main proto type:  [tcp rudp ricmp kcp quic rhttp]")
		f.String("l", "listen", ":8888", "server listen addr")
		f.String("k", "key", "123456", "verify key")
	},
	Run: func(c *grumble.Context) error {
		proto := c.Flags.String("proto")
		listen := c.Flags.String("listen")
		key := c.Flags.String("key")
		if !proxyopt.HasReliableProto(proto) {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "[proto] must be "+fmt.Sprintf("%v", proxyopt.SupportReliableProtos()))
		}
		_, err := net.ResolveTCPAddr("tcp4", listen)
		if err != nil {
			message := fmt.Sprintf("Invalid listen address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		proxyServerParam := proxyopt.ProxyServerParam{
			Proto:  proto,
			Listen: listen,
			Key:    key,
		}
		marshal, _ := json.Marshal(proxyServerParam)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg("Server", common.ClientId, "proxyServerStart", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			if common.GetHttpRetCode(msg) == retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_SUCCESS, common.GetHttpRetData(msg))
			} else {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, common.GetHttpRetData(msg))
			}
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var proxyServerListCmd = &grumble.Command{
	Name:  "list",
	Help:  "list listening proxy servers",
	Usage: "proxy server list [-h | --help]",
	Run: func(c *grumble.Context) error {
		retData := common.SendSuccessMsg("Server", common.ClientId, "proxyServerList", "", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, data := common.SettleRetDataBt(msg)
			var proxyList []proxyopt.ProxyServer
			json.Unmarshal(data, &proxyList)
			proxyopt.PrintServerTable(proxyList)

		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		return nil
	},
}

var proxyServerCloseCmd = &grumble.Command{
	Name:  "close",
	Help:  "turn off the specified proxy listener",
	Usage: "proxy server close [-h | --help] <id>",
	Args: func(a *grumble.Args) {
		a.Int("id", "close proxy id")
	},
	Run: func(c *grumble.Context) error {
		proxyId := c.Args.Int("id")
		var uid string
		retData := common.SendSuccessMsg("Server", common.ClientId, "proxyServerList", "", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		var proxyList []proxyopt.ProxyServer
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, data := common.SettleRetDataBt(msg)
			json.Unmarshal(data, &proxyList)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		closeFlag := false
		for i, proxyServer := range proxyList {
			if proxyId == i+1 {
				uid = proxyServer.Uid
				closeFlag = true
			}
		}
		if !closeFlag {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "The corresponding id does not exist")
			return nil
		}
		data, _ := crypto.Encrypt([]byte(uid), []byte(config.AesKey))
		retData = common.SendSuccessMsg("Server", common.ClientId, "proxyServerClose", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "proxy listener closed successfully")
		return nil
	},
}

var proxyClientCmd = &grumble.Command{
	Name:  "client",
	Help:  "activate the proxy client function",
	Usage: "proxy client [-h | --help]",
}

var proxyClientStartCmd = &grumble.Command{
	Name: "start",
	Help: "connect to proxy server",
	Usage: "proxy client start [-h | --help] [-t | --type type] [-s | --server server] [-p | --proto proto] [-k | --key key] [-F | --fromaddr --fromaddr] [-T | --toaddr]\n" +
		"  eg: \n   proxy client start -s 192.168.1.10:8808 -k admin@123 -F :6000 -T :6000 \n   proxy client start -t reverse_socks5 -s 192.168.1.10:8808 -k admin@123 -F :6666",
	Flags: func(f *grumble.Flags) {
		f.String("t", "type", "bind", "connect type: bind/reverse/bind_socks5/reverse_socks5")
		f.String("s", "server", "127.0.0.1:8888", "server addr")
		f.String("p", "proto", "tcp", "main proto type:  [tcp rudp ricmp kcp quic rhttp]")
		f.String("P", "proxyproto", "tcp", "proxy proto type:  [tcp rudp ricmp kcp quic rhttp]")
		f.String("k", "key", "123456", "verify key")
		f.String("F", "fromaddr", "", "from addr")
		f.String("T", "toaddr", "", "to addr")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		proto := c.Flags.String("proto")
		proxyProto := c.Flags.String("proxyproto")
		serverAddr := c.Flags.String("server")
		connType := c.Flags.String("type")
		key := c.Flags.String("key")
		fromAddr := c.Flags.String("fromaddr")
		toAddr := c.Flags.String("toaddr")
		// 检测参数
		if !proxyopt.HasReliableProto(proto) {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "[proto] must be "+fmt.Sprintf("%v", proxyopt.SupportReliableProtos()))
		}
		if !proxyopt.HasReliableProto(proxyProto) {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "[proxyProto] must be "+fmt.Sprintf("%v", proxyopt.SupportReliableProtos()))
		}
		if !proxyopt.HasReliableConnType(connType) {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "[type] must be "+fmt.Sprintf("%v", proxyopt.SupportReliableConnType()))
		}
		_, err := net.ResolveTCPAddr("tcp4", fromAddr)
		if err != nil {
			message := fmt.Sprintf("Invalid from address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		_, err = net.ResolveTCPAddr("tcp4", toAddr)
		if err != nil {
			message := fmt.Sprintf("Invalid to address [%s]", err)
			colorcode.PrintMessage(colorcode.SIGN_ERROR, message)
			return nil
		}
		proxyClientParam := proxyopt.ProxyClientParam{
			Proto:      proto,
			ProxyProto: proxyProto,
			ConnType:   connType,
			ServerAddr: serverAddr,
			Key:        key,
			FromAddr:   fromAddr,
			ToAddr:     toAddr,
		}
		marshal, _ := json.Marshal(proxyClientParam)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "proxyClientStart", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			fmt.Println(msg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		return nil
	},
}

var proxyClientListCmd = &grumble.Command{
	Name:  "list",
	Help:  "list enabled proxy clients",
	Usage: "proxy client list [-h | --help]",
	Run: func(c *grumble.Context) error {
		retData := common.SendSuccessMsg("Server", common.ClientId, "proxyClientList", "", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		select {
		case msg := <-common.DefaultMsgChan:
			proxyLists := proxyopt.GetClientProxyLists(msg)
			proxyopt.PrintClientTable(proxyLists)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		return nil
	},
}

var proxyClientCloseCmd = &grumble.Command{
	Name:  "close",
	Help:  "turn off the specified proxy client by id",
	Usage: "proxy client close [-h | --help] <id>",
	Args: func(a *grumble.Args) {
		a.Int("id", "close proxy id")
	},
	Run: func(c *grumble.Context) error {
		proxyId := c.Args.Int("id")
		var uid string
		retData := common.SendSuccessMsg("Server", common.ClientId, "proxyClientList", "", "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		var proxyList []proxyopt.ProxyClientList
		select {
		case msg := <-common.DefaultMsgChan:
			data := common.GetHttpRetData(msg)
			json.Unmarshal([]byte(data), &proxyList)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		closeFlag := false
		for i, proxyClient := range proxyList {
			if proxyId == i+1 {
				uid = proxyClient.Uid
				closeFlag = true
			}
		}
		if !closeFlag {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "The corresponding id does not exist")
			return nil
		}
		data, _ := crypto.Encrypt([]byte(uid), []byte(config.AesKey))
		retData = common.SendSuccessMsg(SelectClientId, common.ClientId, "proxyClientClose", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		retData = common.SendSuccessMsg("Server", common.ClientId, "proxyClientClose", data, "")
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "proxy listener closed successfully")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		return nil
	},
}
