package main

import (
	"Orca_Server/routers"
	"Orca_Server/servers"
	"Orca_Server/setting"
	"Orca_Server/sqlmgmt"
	"Orca_Server/tools/log"
	log2 "log"
	"net/http"
)

func init() {
	setting.Setup()
	log.Setup()
}

func main() {
	//打印logo
	//color.Green("OrcaC2 Server " + define.Version)
	//color.Green(define.Logo)
	//初始化数据库
	sqlmgmt.InitDb()
	//初始化路由
	routers.Init()
	//启动一个定时器用来发送心跳
	servers.PingTimer()
	log2.Printf("服务器启动成功，端口号：%s", setting.CommonSetting.HttpPort)

	if err := http.ListenAndServe(":"+setting.CommonSetting.HttpPort, nil); err != nil {
		panic(err)
	}
}
