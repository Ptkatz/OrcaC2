package servers

import (
	"Orca_Server/setting"
	"Orca_Server/tools/crypto"
	"Orca_Server/tools/util"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//channel通道
var ToClientChan chan ClientInfo

//channel通道结构体
type ClientInfo struct {
	ClientId   string
	SendUserId string
	MessageId  string
	Code       int
	Msg        string
	Data       *string
}

type RetData struct {
	MessageId  string      `json:"messageId"`
	SendUserId string      `json:"sendUserId"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

// 心跳间隔
var heartbeatInterval = 25 * time.Second

func init() {
	ToClientChan = make(chan ClientInfo, 1000)
}

var Manager = NewClientManager() // 管理者

func StartWebSocket() {
	websocketHandler := &Controller{}
	http.HandleFunc("/ws", websocketHandler.Run)
	go Manager.Start()
}

//发送信息到指定客户端
func SendMessage2Client(clientId string, sendUserId string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()
	SendMessage2LocalClient(messageId, clientId, sendUserId, code, msg, data)
	return
}

//关闭客户端
func CloseClient(clientId, systemId string) {
	CloseLocalClient(clientId, systemId)
	return
}

//添加客户端到分组
func AddClient2Group(systemId string, groupName string, clientId string, userId string, extend string) {
	if client, err := Manager.GetByClientId(clientId); err == nil {
		//添加到本地group
		Manager.AddClient2LocalGroup(groupName, client, userId, extend)
	}

}

//发送信息到指定分组
func SendMessage2Group(systemId, sendUserId, groupName string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()
	Manager.SendMessage2LocalGroup(systemId, messageId, sendUserId, groupName, code, msg, data)
	return
}

//发送信息到指定系统
func SendMessage2System(systemId, sendUserId string, code int, msg string, data string) {
	encMsg, _ := crypto.Encrypt([]byte(msg), []byte(setting.CommonSetting.CryptoKey))
	messageId := util.GenUUID()
	Manager.SendMessage2LocalSystem(systemId, messageId, sendUserId, code, encMsg, &data)
}

//获取分组列表
func GetOnlineList(systemId *string, groupName *string) map[string]interface{} {
	var clientList []string

	retList := Manager.GetGroupClientList(util.GenGroupKey(*systemId, *groupName))
	clientList = append(clientList, retList...)

	return map[string]interface{}{
		"count": len(clientList),
		"list":  clientList,
	}
}

//通过本服务器发送信息
func SendMessage2LocalClient(messageId, clientId string, sendUserId string, code int, msg string, data *string) {
	//log.WithFields(log.Fields{
	//	"host":     setting.GlobalSetting.LocalHost,
	//	"port":     setting.CommonSetting.HttpPort,
	//	"clientId": clientId,
	//}).Info("发送到通道")
	ToClientChan <- ClientInfo{ClientId: clientId, MessageId: messageId, SendUserId: sendUserId, Code: code, Msg: msg, Data: data}
	return
}

//发送关闭信号
func CloseLocalClient(clientId, systemId string) {
	if conn, err := Manager.GetByClientId(clientId); err == nil && conn != nil {
		if conn.SystemId != systemId {
			return
		}
		Manager.DisConnect <- conn
		log.WithFields(log.Fields{
			"host":     setting.GlobalSetting.LocalHost,
			"port":     setting.CommonSetting.HttpPort,
			"clientId": clientId,
		}).Info("主动踢掉客户端")
	}
	return
}

//监听并发送给客户端信息
func WriteMessage() {
	for {
		clientInfo := <-ToClientChan
		//log.WithFields(log.Fields{
		//	"host":       setting.GlobalSetting.LocalHost,
		//	"port":       setting.CommonSetting.HttpPort,
		//	"clientId":   clientInfo.ClientId,
		//	"messageId":  clientInfo.MessageId,
		//	"sendUserId": clientInfo.SendUserId,
		//	"code":       clientInfo.Code,
		//	"msg":        clientInfo.Msg,
		//	"data":       clientInfo.Data,
		//}).Info("发送到本机")

		if clientInfo.ClientId == "Server" { // 发送给服务器的请求
			msg := clientInfo.Msg
			decMsg, _ := crypto.Decrypt(msg, []byte(setting.CommonSetting.CryptoKey))
			clientId := clientInfo.SendUserId
			data := clientInfo.Data
			SettleMsg(decMsg, clientId, data)
		} else if conn, err := Manager.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {
			if err := Render(conn.Socket, clientInfo.MessageId, clientInfo.SendUserId, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
				Manager.DisConnect <- conn
				log.WithFields(log.Fields{
					"host":     setting.GlobalSetting.LocalHost,
					"port":     setting.CommonSetting.HttpPort,
					"clientId": clientInfo.ClientId,
					"msg":      clientInfo.Msg,
				}).Error("客户端异常离线：" + err.Error())
			}
		}
	}
}

func Render(conn *websocket.Conn, messageId string, sendUserId string, code int, message string, data interface{}) error {
	return conn.WriteJSON(RetData{
		Code:       code,
		MessageId:  messageId,
		SendUserId: sendUserId,
		Msg:        message,
		Data:       data,
	})
}

//启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			for clientId, conn := range Manager.AllClient() {
				if err := conn.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.DisConnect <- conn
					log.Errorf("发送心跳失败: %s 总连接数：%d", clientId, Manager.Count())
				}
			}
		}

	}()
}
