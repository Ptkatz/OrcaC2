package stager

import (
	"Orca_Master/cli/common"
	"Orca_Master/define/api"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"fmt"
	"github.com/togettoyou/wsc"
	"os"
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
		//log.Println("OnConnected: ", Ws.WebSocket.Url)
		// 连接成功后，注册SystemId
		registerSystemId(api.REGISTER_API, config.SystemId)
	})
	common.Ws.OnConnectError(func(err error) {
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "OnConnectError: "+err.Error())
	})
	common.Ws.OnDisconnected(func(err error) {
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "OnDisconnected: "+err.Error())
		os.Exit(10)
	})
	common.Ws.OnClose(func(code int, text string) {
		msg := fmt.Sprintf("OnClose: %d %s", code, text)
		colorcode.PrintMessage(colorcode.SIGN_ERROR, msg)
		done <- true
	})
	common.Ws.OnSentError(func(err error) {
		msg := fmt.Sprintf("OnSentError: %s", err.Error())
		colorcode.PrintMessage(colorcode.SIGN_ERROR, msg)
	})
	common.Ws.OnPingReceived(func(appData string) {
		//fmt.Println("OnPingReceived: ", appData)
		runtime.GC()
	})
	common.Ws.OnPongReceived(func(appData string) {
		//fmt.Println("OnPongReceived: ", appData)
	})
	common.Ws.OnTextMessageSent(func(message string) {
		//fmt.Println("OnTextMessageSent: ", message)
	})
	common.Ws.OnBinaryMessageSent(func(data []byte) {
		//fmt.Println("OnBinaryMessageSent: ", string(data))
	})
	common.Ws.OnTextMessageReceived(func(message string) {
		//fmt.Println("OnTextMessageReceived: ", message)
		SettleMsg(message, common.Ws)
		runtime.GC()
	})
	common.Ws.OnBinaryMessageReceived(func(data []byte) {
		//fmt.Println("OnBinaryMessageReceived: ", string(data))
	})

	// 开始连接
	go common.Ws.Connect()
	for {
		select {
		case <-done:
			return
		}
	}
}
