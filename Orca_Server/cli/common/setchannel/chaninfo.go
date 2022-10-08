package setchannel

import (
	"sync"
)

var mutex sync.Mutex
var CmdMsgChan = make(chan string)                        // 命令消息通道
var FileSliceDataChan = make(map[string]chan interface{}) // 文件元数据通道
