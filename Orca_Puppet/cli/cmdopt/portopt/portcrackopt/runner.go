package portcrackopt

import (
	"Orca_Puppet/cli/cmdopt/portopt/portcrackopt/crack"
	"Orca_Puppet/define/colorcode"
	"fmt"
)

type Runner struct {
	options     *Options
	crackRunner *crack.Runner
}

func NewRunner(options *Options) (*Runner, error) {
	crackOptions := &crack.Options{
		Threads:  options.Threads,
		Timeout:  options.Timeout,
		Delay:    options.Delay,
		CrackAll: options.CrackAll,
		Silent:   options.Silent,
	}
	crackRunner, err := crack.NewRunner(crackOptions)
	if err != nil {
		return nil, fmt.Errorf("crack.NewRunner() err, %v", err)
	}
	return &Runner{
		options:     options,
		crackRunner: crackRunner,
	}, nil
}

func (r *Runner) Run() string {
	var outMessage string
	fmt.Println(r.options.Targets)
	addrs := crack.ParseTargets(r.options.Targets)
	addrs = crack.FilterModule(addrs, r.options.Module)
	if len(addrs) == 0 {
		msg := fmt.Sprintf("target is empty")
		outMessage = colorcode.OutputMessage(colorcode.SIGN_FAIL, msg)
		return outMessage
	}
	// 存活探测
	msg := fmt.Sprintf("detecting live hosts")
	outMessage += colorcode.OutputMessage(colorcode.SIGN_NOTICE, msg)
	addrs = r.crackRunner.CheckAlive(addrs)
	msg = fmt.Sprintf("number of live hosts: %d", len(addrs))
	outMessage += colorcode.OutputMessage(colorcode.SIGN_SUCCESS, msg)
	// 服务爆破
	results, out := r.crackRunner.Run(addrs, r.options.UserDict, r.options.PassDict)
	outMessage += out
	if len(results) > 0 {
		msg = fmt.Sprintf("cracked successfully: %v", len(results))
		outMessage += colorcode.OutputMessage(colorcode.SIGN_SUCCESS, msg)
		for _, result := range results {
			msg = fmt.Sprintf("%v -> %v %v", result.Protocol, result.Addr, result.UserPass)
			outMessage += colorcode.OutputMessage(colorcode.SIGN_SUCCESS, msg)
		}
	}
	return outMessage
}
