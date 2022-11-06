package sqlmgmt

import (
	"Orca_Server/define/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strings"
)

var Db *gorm.DB

func GetDb() *gorm.DB {
	var err error
	pwd, _ := os.Getwd()
	s := strings.Replace(pwd, "\\", "/", -1)
	Db, err := gorm.Open(sqlite.Open(s+config.TeamDB), &gorm.Config{})
	//defer Db.Close()
	if err != nil {
		log.Println("数据库连接失败！", err)
	}
	if err != nil {
		panic(err)
	}
	//Db.LogMode(true)  //sql调试模式
	return Db
}

func InitDb() {
	InitHost()
	InitSshTunnel()
	InitUser()
	InitProxy()
}

func DelRecordByClientId(clientId string) {
	DelProxyRecordByClientId(clientId)
	DelSshTunnelRecordByClientId(clientId)
	DelHostRecordByClientId(clientId)
}

type UsersList struct {
	gorm.Model
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	LoginIp   string `json:"login_ip"`
	LoginTime string `json:"login_time"`
	Online    string `json:"online"`
}

type HostList struct {
	gorm.Model
	ClientId  string `json:"clientId" gorm:"unique"` //唯一标识
	Hostname  string `json:"hostname"`               //主机名
	Ip        string `json:"ip"`                     //上线ip
	ConnPort  string `json:"connPort"`               //上线端口
	Os        string `json:"os"`                     //系统版本
	Privilege string `json:"privilege"`              //客户端执行权限
	Version   string `json:"version"`                //上线客户端版本
	Remarks   string `json:"remarks"`                //备注
}

type AuthList struct {
	gorm.Model
	ClientId string `json:"clientId" gorm:"unique"` //唯一标识
	Token    string `json:"token" gorm:"unique"`    //认证token
	Username string `json:"username"`               // 登录用户
}

type ProxyList struct {
	gorm.Model
	Uid        string `json:"uid" gorm:"unique"`
	ClientId   string `json:"clientId"`
	Proto      string `json:"proto"`
	ProxyProto string `json:"proxyProto"`
	ConnType   string `json:"connType"`
	ServerAddr string `json:"serverAddr"`
	Key        string `json:"key"`
	FromAddr   string `json:"fromAddr"`
	ToAddr     string `json:"toAddr"`
}

type SshtunnelList struct {
	gorm.Model
	Uid      string `json:"uid" gorm:"unique"`
	ClientId string `json:"clientId"`
	Source   string `json:"source"`
	Target   string `json:"target"`
}
