package portcrackopt

import (
	"Orca_Master/tools/util"
	"fmt"
	"strings"
)

type Options struct {
	// input
	Input     string
	InputFile string
	Module    string
	User      string
	Pass      string
	UserFile  string
	PassFile  string
	// config
	Threads  int
	Timeout  int
	Delay    int
	CrackAll bool
	// output
	OutputFile string
	NoColor    bool
	// debug
	Silent bool
	Debug  bool

	Targets  []string
	UserDict []string
	PassDict []string
}

func ParseOptions(input, inputFile, module, user, pass, userFile, passFile string, threads, timeout, delay int, crackAll bool, outputFile string) *Options {
	options := &Options{
		Input:      input,
		InputFile:  inputFile,
		Module:     module,
		User:       user,
		Pass:       pass,
		UserFile:   userFile,
		PassFile:   passFile,
		Threads:    threads,
		Timeout:    timeout,
		Delay:      delay,
		CrackAll:   crackAll,
		OutputFile: outputFile,
		NoColor:    false,
		Silent:     false,
		Debug:      false,
	}
	if err := options.validateOptions(); err != nil {
		fmt.Printf("Program exiting: %v", err)
	}

	if err := options.configureOptions(); err != nil {
		fmt.Printf("Program exiting: %v", err)
	}

	return options
}

// configureOptions 配置选项
func (o *Options) configureOptions() error {
	var err error
	if o.Input != "" {
		o.Targets = append(o.Targets, o.Input)
	} else {
		var lines []string
		lines, err = util.ReadLines(o.InputFile)
		if err != nil {
			return err
		}
		o.Targets = append(o.Targets, lines...)
	}
	if o.User != "" {
		o.UserDict = strings.Split(o.User, ",")
	}
	if o.Pass != "" {
		o.PassDict = strings.Split(o.Pass, ",")
	}
	if o.UserFile != "" {
		if o.UserDict, err = util.ReadLines(o.UserFile); err != nil {
			return err
		}
	}
	if o.PassFile != "" {
		if o.PassDict, err = util.ReadLines(o.PassFile); err != nil {
			return err
		}
	}
	// 去重
	o.Targets = util.RemoveDuplicate(o.Targets)
	o.UserDict = util.RemoveDuplicate(o.UserDict)
	o.PassDict = util.RemoveDuplicate(o.PassDict)

	return nil
}

// validateOptions 验证选项
func (o *Options) validateOptions() error {
	if o.Input == "" && o.InputFile == "" {
		return fmt.Errorf("no service input provided")
	}
	if o.Debug && o.Silent {
		return fmt.Errorf("both debug and silent mode specified")
	}
	if o.Delay < 0 {
		return fmt.Errorf("delay can't be negative")
	}

	return nil
}
