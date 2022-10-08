package routers

import (
	"Orca_Server/api/masterlogin"
	"Orca_Server/api/register"
	"Orca_Server/api/send2client"
	"Orca_Server/servers"
	"net/http"
)

func Init() {
	registerHandler := &register.Controller{}
	masterloginHandler := &masterlogin.Controller{}
	sendToClientHandler := &send2client.Controller{}
	/* 暂时不用的api
	sendToClientsHandler := &send2clients.Controller{}
	sendToGroupHandler := &send2group.Controller{}
	bindToGroupHandler := &bind2group.Controller{}
	getGroupListHandler := &getonlinelist.Controller{}
	closeClientHandler := &closeclient.Controller{}
	*/
	http.HandleFunc("/api/register", registerHandler.Run)
	http.HandleFunc("/api/master_login", masterloginHandler.Run)
	http.HandleFunc("/api/send_to_client", AccessTokenMiddleware(sendToClientHandler.Run))
	//http.HandleFunc("/api/send_to_clients", AccessTokenMiddleware(sendToClientsHandler.Run))
	//http.HandleFunc("/api/send_to_group", AccessTokenMiddleware(sendToGroupHandler.Run))
	//http.HandleFunc("/api/bind_to_group", AccessTokenMiddleware(bindToGroupHandler.Run))
	//http.HandleFunc("/api/get_online_list", AccessTokenMiddleware(getGroupListHandler.Run))
	//http.HandleFunc("/api/close_client", AccessTokenMiddleware(closeClientHandler.Run))

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))

	servers.StartWebSocket()

	go servers.WriteMessage()
}
