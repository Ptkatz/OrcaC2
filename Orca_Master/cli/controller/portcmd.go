package controller

import (
	"Orca_Master/cli/cmdopt/portopt/portcrackopt"
	"Orca_Master/cli/cmdopt/portopt/portscanopt"
	"Orca_Master/cli/cmdopt/sshopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"encoding/json"
	"fmt"
	"github.com/4dogs-cn/TXPortMap/pkg/output"
	"github.com/desertbit/grumble"
	"os"
	"regexp"
	"strings"
	"time"
)

var portCmd = &grumble.Command{
	Name:  "port",
	Help:  "use port scan or port brute",
	Usage: "port scan|crack [-h | --help]",
}

var portScanCmd = &grumble.Command{
	Name: "scan",
	Help: "port Scanner and Banner Identify",
	Usage: "port scan [-h | --help] [-i | --ip ips_or_domain] [-p | --port ports] [--t1000] [-r | --random] [-n | --threads threads_number] [--limit limit_threads_number] [-ei | --eip exclude_ips] [-ep | --eport exclude_ports] [-t | --timeout timeout] [--nbtscan] [-O | --output output_file]\n" +
		"  eg: \n   port scan -i 192.168.1.10 -p 20-6000\n   port scan -i 192.168.1.0/24 --t1000 -r -n 2000 --eport 25,110",
	Flags: func(f *grumble.Flags) {
		f.String("i", "ip", "", "set domain and ips")
		f.String("p", "port", "", "set port ranges to scan，default is top100")
		f.BoolL("t1000", false, "scan top1000 ports")
		f.Bool("r", "random", false, "random scan flag")
		f.Int("n", "threads", 800, "number of threads, between 1 and 2000")
		f.IntL("limit", 0, "lmit number of goroutines, between 1 and 2000")
		f.StringL("eip", "", "set ip ranges to exclude")
		f.StringL("eport", "", "set port ranges to exclude")
		f.Float64("t", "timeout", 0.5, "tcp connect time out default 0.5 second")
		f.BoolL("nbtscan", false, " get netbios stat by UDP137 in local network")
		f.String("O", "output", "", " success log file")
	},
	Run: func(c *grumble.Context) error {
		cmdIps := strings.Split(c.Flags.String("ip"), ",")
		cmdPorts := strings.Split(c.Flags.String("port"), ",")
		cmdT1000 := c.Flags.Bool("t1000")
		cmdRandom := c.Flags.Bool("random")
		numThreads := c.Flags.Int("threads")
		limit := c.Flags.Int("limit")
		excIps := strings.Split(c.Flags.String("eip"), ",")
		excPorts := strings.Split(c.Flags.String("eport"), ",")
		timeout := c.Flags.Float64("timeout")
		nbtscan := c.Flags.Bool("nbtscan")
		outputFile := c.Flags.String("output")
		scanCmdMsg := portscanopt.ScanCmdMsg{
			CmdIps:     cmdIps,
			CmdPorts:   cmdPorts,
			CmdT1000:   cmdT1000,
			CmdRandom:  cmdRandom,
			NumThreads: numThreads,
			Limit:      limit,
			ExcIps:     excIps,
			ExcPorts:   excPorts,
			Tout:       timeout,
			Nbtscan:    nbtscan,
		}
		marshal, _ := json.Marshal(scanCmdMsg)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(sshopt.MySsh.Node, common.ClientId, "portScan", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		var resultEvents []*output.ResultEvent
		portscanopt.Writer, _ = output.NewStandardWriter(false, false, outputFile, "")
		startTime := time.Now()
		colorcode.PrintMessage(colorcode.SIGN_NOTICE, colorcode.COLOR_SHINY+"scanning, please wait..."+colorcode.END)
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, retData := common.SettleRetDataEx(msg)
			_, _, data := common.SettleRetDataBt(msg)
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, strings.Trim(string(data), "\""))
				return nil
			}
			json.Unmarshal(data, &resultEvents)
		case <-time.After(30 * time.Minute):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}
		for _, event := range resultEvents {
			_ = portscanopt.Writer.Write(event)
		}
		if portscanopt.Writer != nil {
			portscanopt.Writer.Close()
		}
		elapsed := time.Since(startTime)
		message := fmt.Sprintf("The scan is over, it takes %v", elapsed)
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, message)
		return nil
	},
}

var portCrackCmd = &grumble.Command{
	Name: "crack",
	Help: "brute-force port cracking, scan weak passwords",
	Usage: "port crack [-h | --help] [-i | --input input] [-f | --input-file input_file] [-m | --module module] [--user user] [--pass pass] [--user-file user_file] [--pass-file pass_file] [-n | --threads threads_number] [-t | --timeout timeout] [--delay delay] [--crack-all crack_all]  [-O | --output output_file]\n" +
		"  eg: \n   port crack -i 192.168.1.10:22 ",
	Flags: func(f *grumble.Flags) {
		f.String("i", "input", "", "crack service input(example: -i '127.0.0.1:3306', -i '127.0.0.1:3307|mysql')")
		f.String("f", "input-file", "", "crack services file(example: -f 'xxx.txt')")
		f.String("m", "module", "all", "choose one module to crack(ftp,ssh,wmi,mssql,oracle,mysql,rdp,postgres,redis,memcached,mongodb)")
		f.StringL("user", "", "user(example: -user 'admin,root')")
		f.StringL("pass", "", "pass(example: -pass 'admin,root')")
		f.StringL("user-file", "", "user file(example: -user-file 'user.txt')")
		f.StringL("pass-file", "", "pass file(example: -pass-file 'pass.txt')")
		f.Int("n", "threads", 1, "number of threads")
		f.Int("t", "timeout", 10, "timeout in seconds")
		f.IntL("delay", 3, "delay between requests in seconds (0 to disable)")
		f.BoolL("crack-all", false, "crack all user:pass")
		f.String("O", "output", "", "output file to write found results")
	},
	Run: func(c *grumble.Context) error {
		input := c.Flags.String("input")
		inputFile := c.Flags.String("input-file")
		module := c.Flags.String("module")
		user := c.Flags.String("user")
		pass := c.Flags.String("pass")
		userFile := c.Flags.String("user-file")
		passFile := c.Flags.String("pass-file")
		threads := c.Flags.Int("threads")
		timeout := c.Flags.Int("timeout")
		delay := c.Flags.Int("delay")
		crackAll := c.Flags.Bool("crack-all")
		outputFile := c.Flags.String("output")
		if threads > 60 {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "the number of threads must not be greater than 60")
			return nil
		}
		options := portcrackopt.ParseOptions(input, inputFile, module, user, pass, userFile, passFile, threads, timeout, delay, crackAll, "")
		marshal, _ := json.Marshal(options)
		data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
		retData := common.SendSuccessMsg(sshopt.MySsh.Node, common.ClientId, "portCrack", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request failed")
			return nil
		}
		colorcode.PrintMessage(colorcode.SIGN_NOTICE, colorcode.COLOR_SHINY+"cracking, please wait..."+colorcode.END)
		select {
		case msg := <-common.DefaultMsgChan:
			_, _, retData := common.SettleRetDataEx(msg)
			_, _, data := common.SettleRetDataBt(msg)
			if retData.Code != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, strings.Trim(string(data), "\""))
				return nil
			}
			fmt.Println(string(data))

			re, _ := regexp.Compile(".*\\[0m ")
			out := re.ReplaceAllString(string(data), "")

			pSaveFile, _ := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR, 0600)
			defer pSaveFile.Close()
			_, err := pSaveFile.Write([]byte(out))
			if err != nil {
				return nil
			}
		case <-time.After(300 * time.Minute):
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "request timed out")
			return nil
		}

		return nil
	},
}
