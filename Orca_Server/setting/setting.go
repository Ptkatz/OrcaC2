package setting

import (
	"Orca_Server/sqlmgmt"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/go-ini/ini"
	"log"
	"net"
	"os"
	"sync"
)

type commonConf struct {
	HttpPort  string
	CryptoKey string
}

var CommonSetting = &commonConf{}

type global struct {
	LocalHost      string //本机内网IP
	ServerList     map[string]string
	ServerListLock sync.RWMutex
}

var GlobalSetting = &global{}

var cfg *ini.File

func Setup() {
	configFile := flag.String("c", "conf/app.ini", "-c conf/app.ini")
	addUserFlag := flag.Bool("au", false, "add user")
	delUserFlag := flag.Bool("du", false, "delete user")
	modUserFlag := flag.Bool("mu", false, "Change the password")
	flag.Parse()
	if *addUserFlag {
		addUser()
		os.Exit(0)
	}

	if *delUserFlag {
		delUser()
		os.Exit(0)
	}

	if *modUserFlag {
		modUserPwd()
		os.Exit(0)
	}

	var err error
	cfg, err = ini.Load(*configFile)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("common", CommonSetting)
	if len(CommonSetting.CryptoKey) > 16 {
		CommonSetting.CryptoKey = CommonSetting.CryptoKey[:16]
	} else if len(CommonSetting.CryptoKey) < 16 {
		CommonSetting.CryptoKey = fmt.Sprintf("%016s", CommonSetting.CryptoKey)
	}

	GlobalSetting = &global{
		LocalHost:  getIntranetIp(),
		ServerList: make(map[string]string),
	}
}

func Default() {
	CommonSetting = &commonConf{
		HttpPort:  "6000",
		CryptoKey: "Adba723b7fe06819",
	}

	GlobalSetting = &global{
		LocalHost:  getIntranetIp(),
		ServerList: make(map[string]string),
	}
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

//获取本机内网IP
func getIntranetIp() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func addUser() {
	username := ""
	password := ""
	repassword := ""
	userprompt := &survey.Input{
		Message: "请输入用户名: ",
	}
	survey.AskOne(userprompt, &username)
	pwdprompt := &survey.Password{
		Message: "请输入密码: ",
	}
	survey.AskOne(pwdprompt, &password)
	pwdprompt = &survey.Password{
		Message: "请再次输入密码: ",
	}
	survey.AskOne(pwdprompt, &repassword)
	if password != repassword {
		log.Fatalf("两次输入的密码不匹配！")
		return
	}
	sqlmgmt.AddUser(username, password)
	log.Println("用户添加成功！")
}

func delUser() {
	usernames := sqlmgmt.GetUsernames()
	selectUsername := ""
	prompt := &survey.Select{
		Message: "请选择要删除的指定用户：",
		Options: usernames,
	}
	survey.AskOne(prompt, &selectUsername)
	sqlmgmt.DelUser(selectUsername)
	log.Println("用户删除成功！")
}

func modUserPwd() {
	usernames := sqlmgmt.GetUsernames()
	selectUsername := ""
	password := ""
	repassword := ""
	prompt := &survey.Select{
		Message: "请选择的指定用户：",
		Options: usernames,
	}
	survey.AskOne(prompt, &selectUsername)
	pwdprompt := &survey.Password{
		Message: "请输入密码: ",
	}
	survey.AskOne(pwdprompt, &password)
	pwdprompt = &survey.Password{
		Message: "请再次输入密码: ",
	}
	survey.AskOne(pwdprompt, &repassword)
	if password != repassword {
		log.Fatalf("两次输入的密码不匹配！")
		return
	}
	sqlmgmt.ModUserPwd(selectUsername, password)
	log.Println("密码修改成功！")
}
