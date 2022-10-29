package crack

import (
	"Orca_Server/cli/cmdopt/portopt/portcrackopt/crack/plugins"
	"Orca_Server/define/colorcode"
	"Orca_Server/tools/util"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"strings"
	"sync"
	"time"
)

type Options struct {
	Threads  int
	Timeout  int
	Delay    int
	CrackAll bool
	Silent   bool
}

type Runner struct {
	options *Options
}

func NewRunner(options *Options) (*Runner, error) {
	return &Runner{
		options: options,
	}, nil
}

type Result struct {
	Addr     string
	Protocol string
	UserPass string
}

type IpAddr struct {
	Ip       string
	Port     int
	Protocol string
}

func (r *Runner) Run(addrs []*IpAddr, userDict []string, passDict []string) (results []*Result, outMessage string) {
	for _, addr := range addrs {
		res, out := r.Crack(addr, userDict, passDict)
		results = append(results, res...)
		outMessage += out
	}
	return
}

func (r *Runner) Crack(addr *IpAddr, userDict []string, passDict []string) (results []*Result, outMessage string) {
	msg := fmt.Sprintf("start to crack: %v:%v %v", addr.Ip, addr.Port, addr.Protocol)
	outMessage += colorcode.OutputMessage(colorcode.SIGN_NOTICE, msg)

	var tasks []plugins.Service
	var taskHash string
	taskHashMap := map[string]bool{}
	// GenTask
	if len(userDict) == 0 {
		userDict = UserMap[addr.Protocol]
	}
	if len(passDict) == 0 {
		passDict = append(passDict, TemplatePass...)
		passDict = append(passDict, CommonPass...)
	}
	for _, user := range userDict {
		for _, pass := range passDict {
			// 替换{user}
			pass = strings.ReplaceAll(pass, "{user}", user)
			// 任务去重
			taskHash = util.Md5(fmt.Sprintf("%v%v%v%v%v", addr.Ip, addr.Port, addr.Protocol, user, pass))
			if taskHashMap[taskHash] {
				continue
			}
			taskHashMap[taskHash] = true
			tasks = append(tasks, plugins.Service{
				Ip:       addr.Ip,
				Port:     addr.Port,
				Protocol: addr.Protocol,
				User:     user,
				Pass:     pass,
				Timeout:  r.options.Timeout,
			})
		}
	}
	// RunTask
	stopHashMap := map[string]bool{}
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	taskChan := make(chan plugins.Service, r.options.Threads)
	for i := 0; i < r.options.Threads; i++ {
		go func() {
			for task := range taskChan {
				addrStr := fmt.Sprintf("%v:%v", addr.Ip, addr.Port)
				userPass := fmt.Sprintf("%v:%v", task.User, task.Pass)
				addrHash := util.Md5(addrStr)
				// 判断是否已经停止爆破
				mutex.Lock()
				if stopHashMap[addrHash] {
					wg.Done()
					mutex.Unlock()
					continue
				}
				mutex.Unlock()
				msg = fmt.Sprintf("[trying] %v", userPass)
				outMessage += colorcode.OutputMessage(colorcode.SIGN_NOTICE, msg)
				scanFunc := plugins.ScanFuncMap[task.Protocol]
				resp := scanFunc(&task)
				switch resp {
				case plugins.CrackSuccess:
					if !r.options.CrackAll {
						mutex.Lock()
						stopHashMap[addrHash] = true
						mutex.Unlock()
					}
					msg = fmt.Sprintf("%v -> %v %v", addr.Protocol, addrStr, userPass)
					outMessage += colorcode.OutputMessage(colorcode.SIGN_NOTICE, msg)
					results = append(results, &Result{
						Addr:     addrStr,
						Protocol: addr.Protocol,
						UserPass: userPass,
					})
				case plugins.CrackError:
					mutex.Lock()
					stopHashMap[addrHash] = true
					mutex.Unlock()
				case plugins.CrackFail:
				}
				if r.options.Delay > 0 {
					time.Sleep(time.Duration(r.options.Delay) * time.Second)
				}
				wg.Done()
			}
		}()
	}

	if r.options.Silent {
		for _, task := range tasks {
			wg.Add(1)
			taskChan <- task
		}
		close(taskChan)
		wg.Wait()
	} else {
		bar := pb.StartNew(len(tasks))
		for _, task := range tasks {
			bar.Increment()
			wg.Add(1)
			taskChan <- task
		}
		close(taskChan)
		wg.Wait()
		bar.Finish()
	}

	msg = fmt.Sprintf("end of crack\n")
	outMessage += colorcode.OutputMessage(colorcode.SIGN_NOTICE, msg)

	return
}
