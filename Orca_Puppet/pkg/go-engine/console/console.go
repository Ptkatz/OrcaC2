package console

import (
	"bufio"
	"fmt"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/loggo"
	"Orca_Puppet/pkg/go-engine/synclist"
	"Orca_Puppet/pkg/go-engine/termcolor"
	"os"
	"strings"
	"sync"
	"time"
)

type Console struct {
	exit           bool
	workResultLock sync.WaitGroup
	readbuffer     chan string
	read           *synclist.List
	write          *synclist.List
	pretext        string
	in             *ConsoleInput
	eb             *EditBox
	normalinput    bool
	color          bool
}

func NewConsole(normalinput bool, historyMaxLen int, color bool) *Console {
	ret := &Console{}
	ret.readbuffer = make(chan string, 16)
	ret.read = synclist.NewList()
	ret.write = synclist.NewList()

	if normalinput {
		go ret.updateRead()
		go ret.run_normal()
	} else {
		ret.in = NewConsoleInput()
		ret.eb = NewEditBox(historyMaxLen)
		err := ret.in.Init()
		if err != nil {
			loggo.Error("NewConsole fail %s", err)
			return nil
		}
		go ret.run()
	}

	return ret
}

func (cc *Console) Stop() {
	cc.exit = true
	cc.workResultLock.Wait()
	if cc.in != nil {
		cc.in.Stop()
	}
}

func (cc *Console) Get() string {
	ret := cc.read.Pop()
	if ret == nil {
		return ""
	}
	return ret.(string)
}

func (cc *Console) Put(str string) {
	cc.write.Push(str)
}

func (cc *Console) Putf(format string, a ...interface{}) {
	cc.write.Push(fmt.Sprintf(format, a...))
}

func (cc *Console) SetPretext(pretext string) {
	cc.pretext = pretext
}

func (cc *Console) Pretext() string {
	return cc.pretext
}

func (cc *Console) updateRead() {
	defer common.CrashLog()

	cc.workResultLock.Add(1)
	defer cc.workResultLock.Done()

	reader := bufio.NewReader(os.Stdin)

	for !cc.exit {
		s, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(time.Millisecond)
			continue
		}
		s = strings.TrimRight(s, "\r")
		s = strings.TrimRight(s, "\n")
		s = strings.TrimSpace(s)
		cc.readbuffer <- s
	}
}

func (cc *Console) run_normal() {
	defer common.CrashLog()

	cc.workResultLock.Add(1)
	defer cc.workResultLock.Done()

	for !cc.exit {
		isneedprintpre := false

		read := ""
		select {
		case c := <-cc.readbuffer:
			read = c
		case <-time.After(time.Duration(100) * time.Millisecond):

		}

		if len(read) > 0 {
			cc.read.Push(read)
			isneedprintpre = true
		}

		for {
			write := cc.write.Pop()
			if write != nil {
				str := write.(string)
				fmt.Println(cc.ouputwithcolor(str))
				isneedprintpre = true
			} else {
				break
			}
		}

		if isneedprintpre {
			fmt.Println(cc.inputwithcolor(cc.pretext + read))
		}
	}
}

func (cc *Console) inputwithcolor(str string) string {
	if cc.color {
		return termcolor.FgString(str, 225, 186, 134)
	} else {
		return str
	}
}

func (cc *Console) ouputwithcolor(str string) string {
	if cc.color {
		return termcolor.FgString(str, 180, 224, 135)
	} else {
		return str
	}
}

func (cc *Console) run() {
	defer common.CrashLog()

	cc.workResultLock.Add(1)
	defer cc.workResultLock.Done()

	for !cc.exit {
		isneedprintpre := false

		for {
			e := cc.in.PollEvent()
			if e == nil {
				break
			}
			isneedprintpre = true
			cc.eb.Input(e)
			str := cc.eb.GetEnterText()
			if len(str) > 0 {
				str = strings.TrimRight(str, "\r")
				str = strings.TrimRight(str, "\n")
				str = strings.TrimSpace(str)
				cc.read.Push(str)
			}
		}

		for {
			write := cc.write.Pop()
			if write != nil {
				str := write.(string)
				fmt.Println(cc.ouputwithcolor(str))
				isneedprintpre = true
			} else {
				break
			}
		}

		if isneedprintpre {
			fmt.Println(cc.inputwithcolor(cc.pretext) + cc.eb.GetShowText(cc.color))
		}
	}
}
